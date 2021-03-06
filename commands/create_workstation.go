package commands

import (
	"os"
	"regexp"

	"github.com/luan/teapot"
	"github.com/luan/tiego/say"
	"github.com/tmtk75/cli"
)

func CreateWorkstation(c *cli.Context) {
	teapotAddr := c.GlobalString("teapot")
	client := teapot.NewClient(teapotAddr)
	name, _ := c.ArgFor("name")
	dockerImage := c.String("docker-image")
	memoryMB := c.Int("memory")
	diskMB := c.Int("disk")
	cpuWeight := uint(c.Int("cpu"))

	err := client.CreateWorkstation(teapot.WorkstationCreateRequest{
		Name: name,
		DockerImage: dockerImage,
		MemoryMB: memoryMB,
		DiskMB: diskMB,
		CPUWeight: cpuWeight,
	})
	if err != nil {
		say.Print(0, say.Bold(say.Red("FAILED: ")))
		var errorMessage string
		if len(err.Error()) > 0 {
			fieldRegexp, _ := regexp.Compile("([^:]*: )([^,]*)(,?)")
			errorMessage = fieldRegexp.ReplaceAllString(err.Error(), "$1"+say.Cyan("$2")+"$3")
		} else {
			errorMessage = "Could not talk to Teapot, did you set the " + say.Cyan("TEAPOT") + " url correctly?"
		}
		say.Println(0, errorMessage)
		os.Exit(1)
	}
	say.Println(0, say.Bold(say.Green("OK")))
}
