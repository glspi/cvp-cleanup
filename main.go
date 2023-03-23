package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"gopkg.in/aristanetworks/go-cvprac.v2/client"
)

// Cleanup Generated Configlets
func cleanup(c *cli.Context) error {
	fmt.Println("Grabbing Configlets")
	cvpClient := c.App.Metadata["client"].(*client.CvpClient)
	reader := bufio.NewReader(os.Stdin)

	lets, err := cvpClient.API.GetConfiglets()
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	for _, let := range lets {
		if let.ContainerCount != 0 && let.NetElementCount != 0 {
			continue
		}
		if let.Type != "Generated" {
			continue
		}
		// Confirm with user before deleting, unless noconfirm set
		if c.Bool("noconfirm") == true {
			cvpClient.API.DeleteConfiglet(let.Name, let.Key)
			fmt.Printf("Deleted %s!\n", let.Name)
		} else {
			for {
				fmt.Printf("Delete %s ? (y/n): ", let.Name)
				response, err := reader.ReadString('\n')
				if err != nil {
					log.Fatal(err)
				}
				response = strings.TrimSpace(response)
				response = strings.ToLower(response)

				if response == "y" || response == "yes" {
					cvpClient.API.DeleteConfiglet(let.Name, let.Key)
					fmt.Printf("Deleted %s!\n", let.Name)
					break
				} else if response == "n" || response == "no" {
					fmt.Println("k then.")
					break
				}
			}
		}
	}

	fmt.Println("Done cleaning generated configlets.")
	return nil
}

// Export ConfigletBuilder
func export(c *cli.Context) error {
	cvpClient := c.App.Metadata["client"].(*client.CvpClient)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Builder Name: ")
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	response = strings.TrimSpace(response)
	builder, err := cvpClient.API.GetConfigletBuilderByName(response)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(builder)
	fout, _ := json.MarshalIndent(builder, "", " ")
	_ = ioutil.WriteFile("test.json", fout, 0644)

	return nil
}

// Login to CVP
func login(c *cli.Context) error {
	ip := c.String("ip")
	user := c.String("user")
	password := c.String("password")

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
	cvpinfo, err := cvpClient.API.GetCvpInfo()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("CVP Version: %s\n", cvpinfo.Version)
	c.App.Metadata["client"] = cvpClient

	return nil
}

// CLI Options
func main() {
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
			&cli.BoolFlag{
				Name:  "noconfirm",
				Usage: "Do not prompt for configlet deletions",
			},
		},
		Commands: []*cli.Command{
			{
				Before: login,
				Name:   "cleanup",
				Usage:  "Cleanup unused configlets in CloudVision",
				Action: cleanup,
			},
			{
				Before: login,
				Name:   "export",
				Usage:  "Export ConfigletBuilder",
				Action: export,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
