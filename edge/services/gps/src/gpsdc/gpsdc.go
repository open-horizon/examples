//
// gpsdc -- a Go client for the gpsd daemon
//
// Usage example:
//
//     package main
//     import (
//       "fmt"
//       "gpsdc"
//     )
//     func main() {
//       var gps_monitor *gpsdc.Session
//       var err error
//       // Open socket to gpsd at specified address (using default here)
//       if gps_monitor, err = gpsdc.Dial(gpsdc.DefaultAddress); err != nil {
//         panic(fmt.Sprintf("Failed to connect to GPSD: ", err))
//       }
//       // Add an inline Filter function
//       gps_monitor.AddFilter("TPV", func (r interface{}) {
//         report := r.(*gpsdc.TPVReport)
//         fmt.Printf("TPV: Lat=%f, Lon=%f, Alt=%f\n",
//                    report.Lat, report.Lon, report.Alt)
//       })
//       // Or a stand-alone Filter function
//       skyfilter := func(r interface{}) {
//         sky := r.(*gpsdc.SKYReport)
//         fmt.Printf("SKY: %d satellites\n", len(sky.Satellites))
//       }
//       gps_monitor.AddFilter("SKY", skyfilter)
//
//       // Tell gpsdc to start watching (using the default watch command here)
//       gps_monitor.SendCommand(gpsdc.DefaultWatchCommand)
//       // Listen forever for reports (delivered to the Filters above)
//       done := gps_monitor.Listen()
//       <- done
//     }
//
// This code is Based on the incomplete and apparently abandoned 2013 project:
//    "https://github.com/stratoberry/go-gpsd"
//
// The majority of this code was copied directly from that project.  It was
// then repaired, and modified to support the more flexible usage pattern
// shown in the example above.
//
// Original contains no copyright or license info
// Original author Josip Lisec, July 9, 2013.
// Modified by Glen Darling, October 2016.
//
package gpsdc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"
	"logutil"
)

// Set logutil.GPS_DEBUG as follows (in logutil.go, or in main.go via env var):
//   0 : No debug output
//   1 : Trace gpsdc interface functions
//   2 : Trace gpsdc internals
//   3 : Dump raw TPV (location) reports
//   4 : Dump all raw reports

const DefaultAddress = "gpsd:2947"
const DefaultWatchCommand = "WATCH={\"enable\":true,\"json\":true}"
const CONNECTION_TIMEOUT=10 	// the timeout when trying to connect to gpsd. It is within the same container so it should either work or fail quickly.

type Filter func(interface{})

type Session struct {
	address string
	socket  net.Conn
	reader  *bufio.Reader
	filters map[string][]Filter
}

type Mode byte

const (
	NoValueSeen Mode = 0
	NoFix       Mode = 1
	Mode2D      Mode = 2
	Mode3D      Mode = 3
)

type GPSDReport struct {
	Class string `json:"class"`
}

type TPVReport struct {
	Class  string    `json:"class"`
	Tag    string    `json:"tag"`
	Device string    `json:"device"`
	Mode   Mode      `json:"mode"`
	Time   time.Time `json:"time"`
	Ept    float64   `json:"ept"`
	Lat    float64   `json:"lat"`
	Lon    float64   `json:"lon"`
	Alt    float64   `json:"alt"`
	Epx    float64   `json:"epx"`
	Epy    float64   `json:"epy"`
	Epv    float64   `json:"epv"`
	Track  float64   `json:"track"`
	Speed  float64   `json:"speed"`
	Climb  float64   `json:"climb"`
	Epd    float64   `json:"epd"`
	Eps    float64   `json:"eps"`
	Epc    float64   `json:"epc"`
}

type SKYReport struct {
	Class      string      `json:"class"`
	Tag        string      `json:"tag"`
	Device     string      `json:"device"`
	Time       time.Time   `json:"time"`
	Xdop       float64     `json:"xdop"`
	Ydop       float64     `json:"ydop"`
	Vdop       float64     `json:"vdop"`
	Tdop       float64     `json:"tdop"`
	Hdop       float64     `json:"hdop"`
	Pdop       float64     `json:"pdop"`
	Gdop       float64     `json:"gdop"`
	Satellites []Satellite `json:"satellites"`
}

type GSTReport struct {
	Class  string    `json:"class"`
	Tag    string    `json:"tag"`
	Device string    `json:"device"`
	Time   time.Time `json:"time"`
	Rms    float64   `json:"rms"`
	Major  float64   `json:"major"`
	Minor  float64   `json:"minor"`
	Orient float64   `json:"orient"`
	Lat    float64   `json:"lat"`
	Lon    float64   `json:"lon"`
	Alt    float64   `json:"alt"`
}

type ATTReport struct {
	Class       string    `json:"class"`
	Tag         string    `json:"tag"`
	Device      string    `json:"device"`
	Time        time.Time `json:"time"`
	Heading     float64   `json:"heading"`
	MagSt       string    `json:"mag_st"`
	Pitch       float64   `json:"pitch"`
	PitchSt     string    `json:"pitch_st"`
	Yaw         float64   `json:"yaw"`
	YawSt       string    `json:"yaw_st"`
	Roll        float64   `json:"roll"`
	RollSt      string    `json:"roll_st"`
	Dip         float64   `json:"dip"`
	MagLen      float64   `json:"mag_len"`
	MagX        float64   `json:"mag_x"`
	MagY        float64   `json:"mag_y"`
	MagZ        float64   `json:"mag_z"`
	AccLen      float64   `json:"acc_len"`
	AccX        float64   `json:"acc_x"`
	AccY        float64   `json:"acc_y"`
	AccZ        float64   `json:"acc_z"`
	GyroX       float64   `json:"gyro_x"`
	GyroY       float64   `json:"gyro_y"`
	Depth       float64   `json:"depth"`
	Temperature float64   `json:"temperature"`
}

type VERSIONReport struct {
	Class      string `json:"class"`
	Release    string `json:"release"`
	Rev        string `json:"rev"`
	ProtoMajor int    `json:"proto_major"`
	ProtoMinor int    `json:"proto_minor"`
	Remote     string `json:"remote"`
}

type DEVICESReport struct {
	Class   string         `json:"class"`
	Devices []DEVICEReport `json:"devices"`
	Remote  string         `json:"remote"`
}

type DEVICEReport struct {
	Class     string  `json:"class"`
	Path      string  `json:"path"`
	Activated string  `json:"activated"`
	Flags     int     `json:"flags"`
	Driver    string  `json:"driver"`
	Subtype   string  `json:"subtype"`
	Bps       int     `json:"bps"`
	Parity    string  `json:"parity"`
	Stopbits  int     `json:"stopbits"`
	Native    int     `json:"native"`
	Cycle     float64 `json:"cycle"`
	Mincycle  float64 `json:"mincycle"`
}

type PPSReport struct {
	Class      string  `json:"class"`
	Device     string  `json:"device"`
	RealSec    float64 `json:"real_sec"`
	RealMusec  float64 `json:"real_musec"`
	ClockSec   float64 `json:"clock_sec"`
	ClockMusec float64 `json:"clock_musec"`
}

type ERRORReport struct {
	Class   string `json:"class"`
	Message string `json:"message"`
}

type Satellite struct {
	PRN  float64 `json:"PRN"`
	Az   float64 `json:"az"`
	El   float64 `json:"el"`
	Ss   float64 `json:"ss"`
	Used bool    `json:"used"`
}

// Dial opens a new connection to the gpsd daemon.
func Dial(addr string) (session *Session, err error) {
	logutil.LogDebug("gpsdc.Dial(%q)\n", addr)
	session = new(Session)
	session.address = addr
	session.filters = make(map[string][]Filter)
	err = session.open()
	return
}

// AddFilter attaches a callback Filter for reports of the specifed class
func (s *Session) AddFilter(class string, f Filter) {
	logutil.LogDebug("gpsdc.AddFilter(%q, <f>)\n", class)
	s.filters[class] = append(s.filters[class], f)
}

// SendCommand sends a command across the socket to the gpsd daemon
func (s *Session) SendCommand(command string) {
	logutil.LogDebug("gpsdc.SendCommand(%q)\n", command)
	fmt.Fprintf(s.socket, "?"+command+";")
}

// Listen runs private listen() function in a goroutine to receive reports
func (s *Session) Listen() (done chan bool) {
	logutil.LogDebug("gpsdc.Listen()\n")
	done = make(chan bool)
	go s.listen(done)
	return
}

// Open a socket to the Session address
func (s *Session) open() (err error) {
	// s.socket, err = net.Dial("tcp4", s.address)  <-- by default, this takes several minutes to time out
	s.socket, err = net.DialTimeout("tcp4", s.address, time.Duration(CONNECTION_TIMEOUT) * time.Second)
	if nil == err {
		s.reader = bufio.NewReader(s.socket)
		if nil == s.reader {
			panic("Unable to allocate new reader.")
		}
	}
	return
}

// Private function to unmarshal a report into an appropriate Go struct
func unmarshal(class string, bytes []byte) (interface{}, error) {
	if logutil.GPS_DEBUG > 1 {
		logutil.LogDebug("* unmarshaling %q\n", class)
	}
	var err error
	switch class {
	case "TPV":
		var r *TPVReport
		err = json.Unmarshal(bytes, &r)
		return r, err
	case "SKY":
		var r *SKYReport
		err = json.Unmarshal(bytes, &r)
		return r, err
	case "GST":
		var r *GSTReport
		err = json.Unmarshal(bytes, &r)
		return r, err
	case "ATT":
		var r *ATTReport
		err = json.Unmarshal(bytes, &r)
		return r, err
	case "VERSION":
		var r *VERSIONReport
		err = json.Unmarshal(bytes, &r)
		return r, err
	case "DEVICES":
		var r *DEVICESReport
		err = json.Unmarshal(bytes, &r)
		return r, err
	case "PPS":
		var r *PPSReport
		err = json.Unmarshal(bytes, &r)
		return r, err
	case "ERROR":
		var r *ERRORReport
		err = json.Unmarshal(bytes, &r)
		return r, err
	}
	return nil, err
}

// Private function to deliver a report to Filters registered for that class
func (s *Session) deliver(class string, report interface{}) {
	for _, f := range s.filters[class] {
		if logutil.GPS_DEBUG > 1 {
			logutil.LogDebug("--> delivering %q to Filter\n", class)
		}
		f(report)
	}
}

// Private function to receive, unmarshal and deliver reports from gpsd daemon
func (s *Session) listen(done chan bool) {
	if logutil.GPS_DEBUG > 1 {
		logutil.Logf("DEBUG: listen()\n")
	}
	for {
		// Read a report line
		if line, err := s.reader.ReadString('\n'); err == nil {
			if logutil.GPS_DEBUG > 3 {
				logutil.LogDebug("<-- %q\n", line)
			}
			// Peek at its class
			var peek GPSDReport
			lineBytes := []byte(line)
			if err = json.Unmarshal(lineBytes, &peek); err == nil {
				if len(peek.Class) == 0 {
					if logutil.GPS_DEBUG > 1 {
						logutil.LogDebug("* no class!\n")
					}
					continue
				}
				if logutil.GPS_DEBUG > 1 {
					logutil.LogDebug("* class=%q\n", peek.Class)
				}
				if logutil.GPS_DEBUG > 2 && "TPV" == peek.Class {
					logutil.LogDebug("* \"TPV\" (raw): %q\n", line)
				}
				// Use the class to unmarshall the data in the report
				if report, err := unmarshal(peek.Class, lineBytes); err == nil {
					// Deliver the report to all registered Filter functions
					s.deliver(peek.Class, report)
				} else {
					logutil.Log("INFO: Ignoring JSON report parsing error:", err)
				}
			} else {
				logutil.Log("INFO: Ignoring JSON class parsing error:", err)
			}
		} else {
			if io.EOF == err {
				for nil != err {
					logutil.Log("WARNING: EOF on stream reader.  Attempting to reconnect...")
					err = s.open()
					time.Sleep(10 * time.Second)
				}
				logutil.Log("INFO: Reconnected to stream reader.")
			} else {
				logutil.Log("INFO: Ignoring stream reader error:", err)
			}
		}
	}
}
