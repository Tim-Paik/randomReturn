/*
Copyright © 2020 Tim_Paik <timpaik@163.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	cfgFile string
	file    string
	path    string
	port    int
	err     error
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "randomReturn",
	Short: "just return something",
	Long:  `just return something`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[randomReturn] Opening server on port " + strconv.Itoa(port))
		fmt.Println("[randomReturn] Listening for " + path)
		gin.DisableConsoleColor()
		gin.SetMode(gin.ReleaseMode)
		var jsons []byte
		var text []string
		var randNumber *rand.Rand
		if jsons, err = ioutil.ReadFile(file); err != nil {
			fmt.Println(err)
			return
		}
		if err := json.Unmarshal(jsons, &text); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("[randomReturn] Reading file " + file)
		r := gin.Default()
		r.GET(path, func(c *gin.Context) {
			c.String(http.StatusOK, "Please enter the keyword after the path!")
		})
		r.GET(path+"/:text", func(c *gin.Context) {
			name := c.Param("text")
			randNumber = rand.New(rand.NewSource(time.Now().UnixNano()))
			returnText := strings.Replace(text[randNumber.Intn(len(text))], "${text}", name, -1)
			c.String(http.StatusOK, returnText+"\n")
		})
		fmt.Println("[randomReturn] Listening for " + path)
		fmt.Println("[randomReturn] ${text} will be replaced with the first path passed in")
		if err := r.Run(":" + strconv.Itoa(port)); err != nil {
			return
		}
		return
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.randomReturn.yaml)")
	rootCmd.PersistentFlags().StringVar(&file, "file", "tmpl.json", "Specifies the replacement string file")
	rootCmd.PersistentFlags().StringVar(&path, "path", "/", "Specifies the path to listen to")
	rootCmd.PersistentFlags().IntVar(&port, "port", 8080, "Specify the port to listen on")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".randomReturn" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".randomReturn")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
