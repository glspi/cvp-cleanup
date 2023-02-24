package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"gopkg.in/aristanetworks/go-cvprac.v2/client"
)

func cleanup(c *cli.Context) error {
	// fmt.Printf("args: %q\n", uhh.Args().Get(0))
	// ip := uhh.String("ip")
	// user := uhh.String("user")
	// password := uhh.String("password")

	// if ip == "" {
	// 	fmt.Println("YOU NEED AN IP ADDRESS!")
	// 	os.Exit(1)
	// }
	// if user == "" {
	// 	fmt.Println("YOU NEED A USERNAME!")
	// 	os.Exit(1)
	// }
	// if password == "" {
	// 	fmt.Println("YOU NEED A PASSWORD!")
	// 	os.Exit(1)
	// }
	// name := "what is this"
	// if uhh.NArg() > 0 {
	// 	name = uhh.Args().Get(1)
	// }
	// if uhh.String("user") == "spanish" {
	// 	fmt.Println("Hola", name)
	// } else {
	// 	fmt.Println("Hello", name)
	// }

	// cvpClient := c.Context.Value("client")
	fmt.Println("Grabbing Configlets") //, c.Context.Value("client").API.GetCvpInfo())
	cvpClient := c.App.Metadata["client"].(*client.CvpClient)
	fmt.Println(cvpClient.API.GetCvpInfo())

	lets, err := cvpClient.API.GetConfigletByName("test1")
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	fmt.Println(lets.Config)

	// data, err := cvpClient.API.GetCvpInfo()
	// if err != nil {
	// 	log.Fatalf("ERROR: %s", err)
	// }
	// fmt.Printf("Data: %v\n", data)

	// bb := uhh.Context.Value("loggedin")
	// fmt.Println("did it again", bb)

	return nil
}

func login(c *cli.Context) error {
	ip := c.String("ip")
	user := c.String("user")
	password := c.String("password")
	fmt.Println(ip, user, password)

	hosts := []string{ip}
	cvpClient, _ := client.NewCvpClient(
		client.Protocol("https"),
		client.Port(443),
		client.Hosts(hosts...),
		client.Debug(false),
	)
	if err := cvpClient.Connect(user, password); err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	fmt.Println("Logged in.")

	c.App.Metadata["client"] = cvpClient

	return nil
}

func main() {
	// lib.Greet()

	//(&cli.App{}).Run(os.Args)
	app := &cli.App{
		Name:  "cvp-cleanup",
		Usage: "CVP Configlet Cleanupper - Usage Instructions:",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "ip",
				Aliases:  []string{"i"},
				Usage:    "IP Address (or fqdn) of CloudVision",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "user",
				Aliases:  []string{"u"},
				Usage:    "Username for CloudVision",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "password",
				Aliases:  []string{"p"},
				Usage:    "Password for CloudVision",
				Required: true,
			},
		},
		Before: login,
		Commands: []*cli.Command{
			{
				Name:   "cleanup",
				Usage:  "Cleanup unused configlets in CloudVision",
				Action: cleanup,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done.")
}
