package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gomodule/redigo/redis"
	_ "github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
	"github.com/rifqiakrm/go-microservice-lib/cache"
	"github.com/rifqiakrm/go-microservice-lib/dialer"
	trace "github.com/rifqiakrm/go-microservice-lib/tracer"
	"github.com/rifqiakrm/go-microservice/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type server interface {
	Run(int) error
}

var (
	cfgFile   string
	dbPool    *sql.DB
	cachePool *redis.Pool
)

var rootCMD = &cobra.Command{
	Use:   "go-microservice",
	Short: "Sample Go Microservice",
	Long:  "Sample Golang microservice",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	cobra.OnInitialize(splash, initconfig, initDB, initCache, GRPCService)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCMD.PersistentFlags().StringVar(&cfgFile, "configs", os.Getenv("PROJECT_ENV"), "configs file (example is $HOME/configs.toml)")
}

// splash print plain text message to console
func splash() {
	fmt.Print(`
                              .__                   .__                                             .__              
  ___________    _____ ______ |  |   ____     _____ |__| ___________  ____  ______ ______________  _|__| ____  ____  
 /  ___/\__  \  /     \\____ \|  | _/ __ \   /     \|  |/ ___\_  __ \/  _ \/  ___// __ \_  __ \  \/ /  |/ ___\/ __ \ 
 \___ \  / __ \|  Y Y  \  |_> >  |_\  ___/  |  Y Y  \  \  \___|  | \(  <_> )___ \\  ___/|  | \/\   /|  \  \__\  ___/ 
/____  >(____  /__|_|  /   __/|____/\___  > |__|_|  /__|\___  >__|   \____/____  >\___  >__|    \_/ |__|\___  >___  >
     \/      \/      \/|__|             \/        \/        \/                 \/     \/                    \/    \/ 
`)
}

func initconfig() {
	viper.SetConfigType("toml")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// search configs in home directory with name "configs" (without extension)
		viper.AddConfigPath("./configs")
		viper.SetConfigName(os.Getenv("CONFIG_FILE"))
	}

	//read env
	viper.AutomaticEnv()

	// if a configs file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("config application:", err)
	}

	log.Println("starting microservice using configs file:", viper.ConfigFileUsed())
}

func initDB() {
	conn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.name"),
		"disable")

	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	dbPool = db

	app.DB = dbPool

	log.Println("database successfully connected!")
}

func Execute() {
	if err := rootCMD.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initCache() {
	redisHost := fmt.Sprintf("%s:%s", viper.GetString("redis.host"), viper.GetString("redis.port"))
	if redisHost == ":" {
		redisHost = "localhost:6379"
	}
	cachePool = newPool(redisHost)

	ctx := context.Background()
	_, err := cachePool.GetContext(ctx)

	if err != nil {
		panic("failed to connect to redis")
	}

	cache.Init(cachePool)

	log.Println("redis successfully connected!")
	cleanupHook()
}

func newPool(server string) *redis.Pool {
	return &redis.Pool{

		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func cleanupHook() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		cachePool.Close()
		os.Exit(0)
	}()
}

func GRPCService() {
	var port = viper.GetInt("app.port")
	var jaegeraddr = fmt.Sprintf("%v:%v",
		viper.GetString("jaeger.host"),
		viper.GetString("jaeger.port"),
	)
	tracer, err := trace.New(viper.GetString("app.name"), jaegeraddr)
	if err != nil {
		log.Fatalf("trace new error: %v", err)
	}
	log.Println("jaeger initiated!")

	var srv server

	srv = app.NewSample(tracer,
		initGRPCConn(fmt.Sprintf("%v:%v", viper.GetString("another_rpc.host"), viper.GetString("another_rpc.port")), tracer),
	)

	if err := srv.Run(port); err != nil {
		log.Fatalf("failed to start rpc server : %v", err)
	}
}

func initGRPCConn(addr string, tracer opentracing.Tracer) *grpc.ClientConn {
	conn, err := dialer.Dial(addr, dialer.WithTracer(tracer))
	if err != nil {
		panic(fmt.Sprintf("ERROR: dial error: %v", err))
	}
	return conn
}
