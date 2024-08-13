package video

import (
	"log"
	"time"
)

type Verb uint16

const (
	GET Verb = iota
	SET
	HIDEALL
	RECORD_START
	RECORD_STOP
)

const (
	RecordingFolder = "recordings/"
)

var cmdList = []string{
	"Get",
	"Set",
	"HideAll",
}

func (cmd Verb) String() string {
	if cmd >= Verb(len(cmdList)) {
		return "Unknown"
	}
	return cmdList[cmd]
}

type ServerCmd struct {
	Action Verb
	Value  any
}

type Server struct {
	Source VideoSource
	Config *VideoConfig

	Quit chan int
	Cmd  chan ServerCmd

	StreamHook *StreamHook

	Filters []Hook

	HideMain  bool
	HideThumb bool
	HideAll   bool
	Busy      bool

	Recording  bool
	recordStop time.Time
	buffer     []byte

	captureCount int64
	captureTotal int64
	captureEnd   time.Time

	captureStop   chan int
	captureSource chan []byte
}

func NewVideoServer(source VideoSource, config *VideoConfig) *Server {

	cam := &Server{
		Source:        source,
		Config:        config,
		Quit:          make(chan int),
		Cmd:           make(chan ServerCmd),
		StreamHook:    NewStreamHook(),
		Filters:       make([]Hook, 0),
		captureStop:   make(chan int),
		captureSource: make(chan []byte),
	}

	return cam
}

func (vs *Server) GetSource() VideoSource {
	return vs.Source
}

func (vs *Server) AddFilter(filter Hook) {
	vs.Filters = append(vs.Filters, filter)
}
func (vs *Server) Command(cmd ServerCmd) {
	vs.Cmd <- cmd
}

func (vs *Server) RecordCmd(seconds int) {
	vs.Command(ServerCmd{Action: RECORD_START, Value: seconds})
}

func (vs *Server) StopRecordCmd() {
	vs.Command(ServerCmd{Action: RECORD_STOP, Value: true})
}

func (vs *Server) Open() (err error) {
	err = vs.Source.Open(vs.Config)
	if err != nil {
		log.Printf("Open Error '%s', %v\n", vs.Source.Path(), err)
	} else {
		log.Printf("Opened '%s'\n", vs.Source.Path())
	}
	return
}

func (vs *Server) Close() {
	if vs.Recording {
		vs.stopRecording()
	}
	vs.Source.Close()
	log.Printf("Closed '%s'\n", vs.Source.Path())
}

const (
	delayNormal    = time.Millisecond
	delayRetry     = time.Second
	delayHibernate = time.Second * 30
	recordLimit    = time.Second * 5
)

func (vs *Server) startRecording(duration int) {
	log.Println("start recording")

	if vs.Recording {
		log.Println("already recording")
		vs.stopRecording()
		return //?
	}

	vs.Recording = true
	vs.captureCount = 0
	config := vs.Source.Config()

	go Capture(vs.captureStop, vs.captureSource,
		config.Width, config.Height, config.FPS)

	now := time.Now()
	vs.recordStop = now.Add(time.Second * time.Duration(duration))
	log.Println("recording started...")

}

func (cam *Server) stopRecording() {
	if !cam.Recording {
		log.Println("stopRecording already stopped")
		return
	}

	cam.captureStop <- 1
	cam.Recording = false
	log.Println("recorder closed")
}

func (vs *Server) doCmd(cmd ServerCmd) {
	switch cmd.Action {
	// case GET:
	// 	cmd.Value = cam.video.Get(cmd.Property)
	// case SET:
	// 	f, _ := cmd.Value.(float64)
	// 	cam.video.Set(cmd.Property, float64(f))
	case HIDEALL:
		b, _ := cmd.Value.(bool)
		vs.HideAll = b
	case RECORD_START:
		vs.startRecording(cmd.Value.(int))
	case RECORD_STOP:
		vs.stopRecording()
	}
}

func (vs *Server) Serve() {
	if vs.Busy {
		log.Fatal("server already busy")
		return
	}

	if !vs.Source.IsOpened() {
		log.Fatal("source not opened")
		return
	}

	vs.Busy = true
	defer func() {
		vs.Busy = false
		vs.Close()
	}()

	var (
		cmd   ServerCmd
		retry int = 0
	)

	var (
		delay = delayNormal
		buf   []byte
		err   error
	)

	// cam.buffer

	for {
		time.Sleep(delay)

		select {
		case <-vs.Quit:
			return
		case cmd = <-vs.Cmd:
			vs.doCmd(cmd)
			continue
		default:
		}

		if vs.HideAll {
			continue
		}

		buf, err = vs.Source.Read()
		if err != nil {
			log.Println("READ", err)
			if retry > 10 {
				delay = delayHibernate
			} else {
				delay = delayRetry
			}

			retry++
			log.Printf("%v is unavailable, attempts=%d next in %.0f seconds\n",
				vs.Source.Path(), retry, delay.Seconds())
			if vs.Source.IsOpened() {
				vs.Source.Close()
			}
			vs.Open()
			continue
		}
		delay = delayNormal
		retry = 0
		vs.StreamHook.Update(buf)

		// fmt.Println(len(buf), "BYTES READ")

		// if cam.matrix.Empty() {
		// 	continue
		// }

		// for _, filter := range vs.Filters {
		// 	filter.Update(&vs.buffer)
		// }

		if vs.Recording {
			// log.Println("recording")
			vs.captureSource <- buf
			if vs.recordStop.Before(time.Now()) {
				vs.stopRecording()
			}
		}
	}

}
