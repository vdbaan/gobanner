// Copyright © 2019 S. van der Baan <steven@vdbaan.net>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// bannergrab in Go
// based on bannergrab 3.5
package main


import (
    // "fmt"
    "os"
    "time"
    "strconv"
    // "errors"
    "github.com/op/go-logging"
    "github.com/spf13/cobra"
    "github.com/fatih/color"
    "net"
    // "crypto/tls"
)

var (
    rootCmd = &cobra.Command{
        Use:   "gobanner [flags] HOST PORT",
        Version: programVersion,
        Short: "GoBanner - a network service banner grabbing tool.",
        Long: 
`GoBanner performs connection, trigger and basic service information collection. 
There are basic banner grabbing modes, the first mode (the default one) sends 
triggers to the services and performs basic information collection. The second 
mode (--no-triggers), only connects to the service and returns the connection 
banner.`,

        Run: func(cmd *cobra.Command, args []string) {
            color.Magenta(programBanner)
            if verbose {
                backendLeveled.SetLevel(logging.DEBUG, "")
            } else {
                backendLeveled.SetLevel(logging.INFO, "")
            }
            log.Debug("Verbose: ON")
            if showTriggers {
                ShowTriggers()
                os.Exit(0)
            }
            if version {
                log.Info(programVersion)
                os.Exit(0)
            }
            if len(args) < 2 {
                cmd.Help()
                os.Exit(0)
            }
            host = args[0]
            port,_ = strconv.Atoi(args[1])
            goBannerGrab()
        },
	}

	log            = logging.MustGetLogger("gobanner")
	logFormat      = logging.MustStringFormatter(`%{color}%{message}%{color:reset}`)
	backendLeveled logging.LeveledBackend

    host string
    port int
	udp bool
	noTriggers bool
	trigger string
	noSsl bool
	noHex bool
	verbose bool
	showTriggers bool
    version bool
    connTime int
    readTime int 
)

// main sets the logging and executes the cobra command
func main() {
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, logFormat)
	backendLeveled = logging.AddModuleLevel(backendFormatter)
    logging.SetBackend(backendLeveled)
        
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// init specifies the flags that are supported
func init() {
    rootCmd.PersistentFlags().BoolVarP(&udp, "udp","",false,"Connect to a port using UDP. The default is to use TCP.")
    rootCmd.PersistentFlags().BoolVarP(&noTriggers, "no-triggers","",false,"Collect only the connection banner,  no triggers and no SSL.")
    rootCmd.PersistentFlags().StringVarP(&trigger,"trigger","","","Specify  the  trigger  to use.  Specify DEFAULT to use the default trigger.")
    rootCmd.PersistentFlags().BoolVarP(&noSsl,"no-ssl","",false,"Prevent SSL connection creation.")
    rootCmd.PersistentFlags().BoolVarP(&noHex, "no-hex","",false,"Output containing non-printable characters are converted to hex. This option prevents the conversion.")
    rootCmd.PersistentFlags().IntVarP(&connTime,"conn-time","",5,"Connection timeout in seconds(default is 5s).")
    rootCmd.PersistentFlags().IntVarP(&readTime, "read-time","",3,"Read timeout in seconds (default is 3s).")
    rootCmd.PersistentFlags().BoolVarP(&verbose,"verbose","",false,"Show additional program details such as any errors.")
    rootCmd.PersistentFlags().BoolVarP(&showTriggers,"show-triggers","",false,"Show the supported triggers.")
    rootCmd.Flags().SortFlags = false
    rootCmd.PersistentFlags().SortFlags = false
}


// goBannerGrab tries to grab the banner of a service
// it will create a connection to the server and either tries to retrieve the default response if no triggers are selected
// otherwise it will try to select the correct method (either forced by a selected trigger, or based on port number) and
// send a trigger to elicit a response.
func goBannerGrab() {
    conn, err := getConnection()
    ifErrorMessageStop(err, "Problem connecting to client")    
    if noTriggers {
        getDefaultBanner(conn)
    } else {
        var service serviceConfig 
        if trigger == "" {
            // get trigger based on port
            for _,t := range services {
                if t.port == port {
                    service = t
                    log.Debugf("Using trigger based on port %d",port)
                }
            }
        } else if trigger == "DEFAULT" {
            service = services["DEFAULT"]
        } else {
            service = services[trigger]
        }
        if service.trigger == nil {
            log.Error("Can't find requested service")
        }
        var len int
        reply := make([] byte ,1024)

        if service.connectBanner {
            getDefaultBanner(conn)
        }
        for i,t := range service.trigger {
            log.Debugf("Sending trigger #%d\n%s",i, printOutputS(t.trigger))
            _, err = conn.Write([]byte(t.trigger))
            ifErrorMessageStop(err, "Error writing to server")
            log.Debugf("Reading reply of trigger #%d",i)
            len, err = conn.Read(reply)
            ifErrorMessageStop(err, "Error reading from server")
            log.Infof("received:\n%s",printOutput(reply, len))
        }       
    }
     conn.Close()
}

// getDefaultBanner simply reads and prints from the provided connection
func getDefaultBanner(conn net.Conn)  {
    reply := make([] byte ,1024)
    len, err := conn.Read(reply)
    ifErrorMessageStop(err, "Error reading from server")
    log.Info("received: %s",printOutput(reply,len))
}

// getConnection creates a connection based on the provided flags
func getConnection() (net.Conn, error) {
    var protocol string
    if udp {
        protocol = "udp"
    } else {
        protocol = "tcp"
    }
    log.Debugf("Using %s protocol",protocol)

    log.Debugf("Setting Connection timeout to %d seconds", connTime)
    conn, err := net.DialTimeout(protocol, host+":"+strconv.Itoa(port),time.Duration(connTime) +time.Second)

    log.Debugf("Setting ReadDeadline to %d seconds", readTime)
    conn.SetReadDeadline(time.Now().Add(time.Duration(readTime) + time.Second))
    return conn,err
    // -- SSL
    // conf := &tls.Config{
         //InsecureSkipVerify: true,
    //}
    // conn, err := tls.Dial("tcp", "127.0.0.1:443", conf)
}