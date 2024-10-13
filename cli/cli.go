package cli

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"grpc-todo-list/config"
	"grpc-todo-list/db"
	"grpc-todo-list/pb"
	"grpc-todo-list/service"
)

// Global variable to hold the config
var c config.Config

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "Todo CLI",
	Run: func(cmd *cobra.Command, args []string) {
		// Connect to the database and start the gRPC server
		if err := startServer(c); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	},
}

func Execute(config config.Config) error {
	c = config
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func init() {
	// Add flags for database connection and gRPC port
	rootCmd.Flags().StringVar(&c.DB.User, "db-user", "root", "Database username")
	rootCmd.Flags().StringVar(&c.DB.Password, "db-password", "my_secure_password", "Database password")
	rootCmd.Flags().StringVar(&c.DB.Host, "db-host", "db", "Database host")
	rootCmd.Flags().StringVar(&c.DB.Port, "db-port", "3306", "Database port")
	rootCmd.Flags().StringVar(&c.DB.Name, "db-name", "todo_app", "Database name")
	rootCmd.Flags().StringVar(&c.GrpcPort, "grpc-port", "50054", "gRPC server port")
	rootCmd.Flags().StringVar(&c.HttpPort, "http-port", "4000", "HTTP server port")
}

func startServer(conf config.Config) error {
	// Create a parent context and a cancel function for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a listener on the specified gRPC port
	listener, err := net.Listen("tcp", ":"+conf.GrpcPort)
	if err != nil {
		return err
	}

	// Set up signal handling for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		fmt.Println("Received shutdown signal, shutting down...")
		cancel() // Call the cancel function to exit the context
	}()

	// Initialize the FX application
	fx.New(
		fx.Provide(
			func() (*db.Database, error) {
				// Construct the DSN from the flags
				dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
					conf.DB.User, conf.DB.Password, conf.DB.Host, conf.DB.Port, conf.DB.Name)
				return db.NewDatabase(dsn)
			},
			service.NewTodoService,
		),
		fx.Invoke(func(todoService *service.TodoService) {
			grpcServer := grpc.NewServer()
			pb.RegisterTodoServiceServer(grpcServer, todoService)

			go func() {
				if err := grpcServer.Serve(listener); err != nil {
					log.Fatalf("failed to serve: %v", err)
				}
			}()

			// Start the HTTP server for gRPC-Gateway
			go startHTTPServer(ctx)

			// Start the interactive client command loop with graceful shutdown
			go interactiveClient(ctx) // Pass context to interactive client
		}),
	).Run()

	return nil
}

func startHTTPServer(ctx context.Context) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterTodoServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%s", c.GrpcPort), opts)
	if err != nil {
		log.Fatalf("Failed to register gRPC Gateway: %v", err)
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", c.HttpPort),
		Handler: mux,
	}

	go func() {
		fmt.Printf("HTTP server listening on :%s\n", c.HttpPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	fmt.Println("Shutting down HTTP server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server shutdown failed: %v", err)
	}
}

func interactiveClient(ctx context.Context) {
	reader := bufio.NewReader(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Exiting interactive client...")
			return // Exit the loop if the context is done
		default:
			fmt.Println("\nAvailable commands:")
			fmt.Println("[1] Add a Task")
			fmt.Println("[2] Get all Tasks")
			fmt.Println("[3] Complete a Task")
			fmt.Println("[4] Exit")
			fmt.Print("\nChoose a command number: ")

			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			switch input {
			case "1":
				addTask(reader)
			case "2":
				getTasks()
			case "3":
				completeTask(reader)
			case "4":
				fmt.Println("Exiting interactive client...")
				return
			default:
				fmt.Println("Invalid choice, please enter a valid command number.")
			}
		}
	}
}

func addTask(reader *bufio.Reader) {
	fmt.Print("\nEnter Task Title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	fmt.Print("Enter Task Description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%s", c.GrpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTodoServiceClient(conn)

	req := &pb.AddTaskRequest{
		Title:       title,
		Description: description,
	}

	_, err = client.AddTask(context.Background(), req)
	if err != nil {
		log.Fatalf("could not add task: %v", err)
	}

	fmt.Println("Task added successfully!")
}

func completeTask(reader *bufio.Reader) {
	fmt.Print("\nEnter Task ID to Complete: ")
	taskID, _ := reader.ReadString('\n')
	taskID = strings.TrimSpace(taskID)

	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%s", c.GrpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTodoServiceClient(conn)

	var id int32
	_, err = fmt.Sscan(taskID, &id)
	if err != nil {
		fmt.Println("")
	}

	req := &pb.CompleteTaskRequest{
		Id: id,
	}

	_, err = client.CompleteTask(context.Background(), req)
	if err != nil {
		log.Fatalf("could not complete task: %v", err)
	}

	fmt.Println("Task completed successfully!")
}

func getTasks() {
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%s", c.GrpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTodoServiceClient(conn)

	resp, err := client.GetTasks(context.Background(), &pb.GetTasksRequest{})
	if err != nil {
		log.Fatalf("could not get tasks: %v", err)
	}

	fmt.Println("Tasks: ")
	for _, task := range resp.Tasks {
		status := "Incomplete"
		if task.IsCompleted {
			status = "Completed"
		}

		fmt.Println(task.CreatedAt)
		fmt.Printf("ID: %d, Title: %s, Description: %s, Status: %s\n", task.Id, task.Title, task.Description, status)
	}
}
