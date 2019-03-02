package godash

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
)

var (
	DASH_SERVICE = ble.MustParse("af237777879d61861f49deca0e85d9c1")
	DASH_NOTIFY  = ble.MustParse("af230006879d61861f49deca0e85d9c1")

	DASH_CMD = ble.MustParse("af230002879d61861f49deca0e85d9c1")
)

type Dash struct {
	device ble.Device
	client ble.Client

	dashCmd *ble.Characteristic

	connectedCh chan error
	err         error

	// sensor data
	ProxRight byte
	ProxLeft  byte
	ProxRear  byte
}

func New() *Dash {
	return &Dash{
		connectedCh: make(chan error),
	}
}

func (d *Dash) Err() error {
	return d.err
}

func (d *Dash) Connect() error {
	dev, err := linux.NewDevice()
	if err != nil {
		return err
	}
	d.device = dev
	ble.SetDefaultDevice(d.device)

	d.client, err = ble.Connect(context.TODO(), func(a ble.Advertisement) bool {
		return a.LocalName() == "Dash"
	})
	if err != nil {
		return err
	}
	p, err := d.client.DiscoverProfile(true)
	if err != nil {
		return fmt.Errorf("cannot discover descriptors: %s", err)
	}
	char := p.FindCharacteristic(ble.NewCharacteristic(DASH_NOTIFY))
	if err := d.client.Subscribe(char, false, d.decodeDashNotify); err != nil {
		return fmt.Errorf("cannot subscribe: %s", err)
	}
	d.dashCmd = p.FindCharacteristic(ble.NewCharacteristic(DASH_CMD))

	return nil
}

func (d *Dash) Disconnect() {
	d.client.CancelConnection()
}

func (d *Dash) decodeDashNotify(b []byte) {
	d.ProxLeft = b[6]
	d.ProxRight = b[7]
	d.ProxRear = b[8]
}

func (d *Dash) Stop() error {
	return d.command("drive", []byte{0, 0, 0})
}

func (d *Dash) Drive(speed int) error {
	return d.command("drive", []byte{byte(speed & 0xff), 0x00, byte((speed & 0x0f00) >> 8)})
}

func (d *Dash) Eye(v int16) error {
	return d.command("eye", []byte{uint8(v >> 8), uint8(v & 0xff)})
}

func (d *Dash) Turn(degrees int) error {
	// FIXME: extract into proper helper
	eight_byte := 0x80

	distance_mm := 0
	seconds := math.Abs((float64(degrees) / (360.0 / 2.094)))

	sixth_byte := 0
	seventh_byte := 0

	distance_low_byte := distance_mm & 0x00ff
	distance_high_byte := (distance_mm & 0x3f00) >> 8
	sixth_byte |= distance_high_byte

	centiradians := int(float64(degrees) * math.Pi / 180.0 * 100.0)
	turn_low_byte := centiradians & 0x00ff
	turn_high_byte := (centiradians & 0x0300) >> 2
	sixth_byte |= turn_high_byte
	if centiradians < 0 {
		seventh_byte = 0xc0
	}
	time_measure := int(seconds * 1000.0)
	time_low_byte := time_measure & 0x00ff
	time_high_byte := (time_measure & 0xff00) >> 8

	b := []byte{
		byte(distance_low_byte),
		byte(0x00), //unknown
		byte(turn_low_byte),
		byte(time_high_byte),
		byte(time_low_byte),
		byte(sixth_byte),
		byte(seventh_byte),
		byte(eight_byte),
	}
	err := d.command("move", b)
	time.Sleep(time.Duration(seconds * float64(time.Second)))
	return err
}

func (d *Dash) command(name string, values []byte) error {
	cmds := map[string]byte{
		"drive": 0x02,
		"eye":   0x09,
		"move":  0x23,
	}
	b := bytes.NewBuffer(nil)
	b.Write([]byte{cmds[name]})
	b.Write(values)
	return d.client.WriteCharacteristic(d.dashCmd, b.Bytes(), true)
}
