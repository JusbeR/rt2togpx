package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/tkrajina/gpxgo/gpx"
)

type RT2Waypoint struct {
	SomeField    string
	WaypointName string
	Latitude     float64
	Longitude    float64
}

type RT2 struct {
	H1        string
	H2        string
	H3        string
	Waypoints []RT2Waypoint
}

// H1,Maphelper VMP(Virtual Map Points) File Version 1.0
// H2,WGS 84
// H3,voikoski2museo,,0
// W,RTP1,61.254227,26.758960,0
// W,RTP2,61.253320,26.759020,1

func parseRT2Line(line string) (rT2Waypoint RT2Waypoint, err error) {
	slice := strings.Split(line, ",")
	if len(slice) < 4 {
		return rT2Waypoint, fmt.Errorf("Invalid line '%s'", line)
	}
	rT2Waypoint.SomeField = slice[0]
	rT2Waypoint.WaypointName = slice[1]
	rT2Waypoint.Latitude, err = strconv.ParseFloat(slice[2], 64)
	if err != nil {
		return rT2Waypoint, fmt.Errorf("Failed to parse latitude from line '%s'", line)
	}
	rT2Waypoint.Longitude, err = strconv.ParseFloat(slice[3], 64)
	if err != nil {
		return rT2Waypoint, fmt.Errorf("Failed to parse longitude from line '%s'", line)
	}
	return rT2Waypoint, nil
}

func readRT2File(fileName string) (rt2 RT2, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return rt2, fmt.Errorf("Failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "H1,") {
			rt2.H1 = line
		} else if strings.Contains(line, "H2,") {
			rt2.H2 = line
		} else if strings.Contains(line, "H3,") {
			rt2.H3 = line
		} else if strings.Contains(line, "W,") {
			rT2Waypoint, err := parseRT2Line(line)
			if err != nil {
				return rt2, err
			}
			rt2.Waypoints = append(rt2.Waypoints, rT2Waypoint)
		}
	}

	if err := scanner.Err(); err != nil {
		return rt2, fmt.Errorf("Failed read file: %v", err)
	}
	return rt2, nil
}

func gpxPointsFromRt2(rt2 RT2) (gpxPoints []gpx.GPXPoint) {
	for _, point := range rt2.Waypoints {
		point := gpx.Point{
			Latitude:  point.Latitude,
			Longitude: point.Longitude,
		}
		gpxPoint := gpx.GPXPoint{
			Point: point,
		}
		gpxPoints = append(gpxPoints, gpxPoint)
	}
	return gpxPoints
}

// TODO: Read all the info from RT2 and try to push that to
// correct place in gpx
func main() {
	fileIn := flag.String("rt2file", "", "rt2 file to be converted")
	fileOut := flag.String("out", "route.gpx", "output gpx file name")
	verbose := flag.Bool("verbose", false, "Print some details about route")
	flag.Parse()

	if len(*fileIn) < 1 {
		flag.PrintDefaults()
		log.Fatalln("You must give 'rt2file' to convert")
	}
	rt2, err := readRT2File(*fileIn)
	if err != nil {
		log.Fatalln("Failed to parse RT2 file:", err)
	}
	gpxPoints := gpxPointsFromRt2(rt2)

	gpxFile := gpx.GPX{}
	gpxFile.Creator = "https://github.com/JusbeR/rt2togpx"
	gpxRoute := gpx.GPXRoute{}
	gpxFile.Routes = append(gpxFile.Routes, gpxRoute)
	gpxFile.Routes[0].Points = gpxPoints

	xmlBytes, err := gpxFile.ToXml(gpx.ToXmlParams{Version: "1.1", Indent: true})
	if err != nil {
		log.Fatalln("Failed to form gpx file:", err)
	}
	err = ioutil.WriteFile(*fileOut, xmlBytes, 0644)
	if err != nil {
		log.Fatalln("Failed to write gpx file:", err)
	}
	if *verbose {
		for _, route := range gpxFile.Routes {
			log.Printf("Route length: %.2fkm", route.Length()/1000)
		}
	}
	log.Println("GPX file succesfully saved to", *fileOut)
}
