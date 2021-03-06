package nmea

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"
	"time"
	"unicode"
)

func TestChecksum(t *testing.T) {
	tests := map[string]bool{
		"":     false,
		"*00":  false,
		"$*00": true,
		"$*01": false,
		"^0*0": false,
		"$0*0": false,
		"$*xx": false,
		"$GPRMC,162254.00,A,3723.02837,N,12159.39853,W,0.820,188.36,110706,,,A*74": true,
		"$GPRMC,162254.00,A,3723.02837,N,12159.39853,W,0.820,188.36,110706,,,A*72": false,
	}

	for in, exp := range tests {
		if checkChecksum(in) != exp {
			t.Errorf("Failed on %v/%v", in, exp)
		}
	}
}

func TestSampleParsing(t *testing.T) {
	for _, s := range strings.Split(ubloxSample, "\n") {
		parseMessage(s, nil)
	}
}

func TestSampleProcessing(t *testing.T) {
	err := Process(strings.NewReader(ubloxSample), nil, nil)
	if err != nil {
		t.Errorf("Unexpected error, got %v", err)
	}
}

type rmcHandler struct {
	rmc RMC
}

func (r *rmcHandler) HandleRMC(rmc RMC) {
	r.rmc = rmc
}

func logJSON(t *testing.T, h interface{}) {
	j, err := json.Marshal(h)
	if err != nil {
		t.Errorf("Failed to marshal %v: %v", h, err)
	}
	t.Logf("%T: %s", h, j)
}

const ε = 0.00001

func near(a, b float64) bool {
	return math.Abs(a-b) < ε
}

func similar(t *testing.T, a, b interface{}) bool {
	ta := reflect.TypeOf(a)
	tb := reflect.TypeOf(b)
	if ta != tb {
		t.Errorf("Expected same type between %v and %v", ta, tb)
		return false
	}
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	for i := 0; i < ta.NumField(); i++ {
		f := ta.Field(i)
		name := f.Name
		if !unicode.IsUpper(rune(name[0])) {
			continue
		}
		af := va.Field(i)
		bf := vb.Field(i)
		if af.Type() != bf.Type() {
			t.Errorf("Incorrect type in field %v: %T != %T", name, af.Type(), bf.Type())
			return false
		}
		av := af.Interface()
		bv := bf.Interface()

		switch av.(type) {
		case time.Time:
			if !av.(time.Time).Equal(bv.(time.Time)) {
				t.Errorf("Timestamp field %v was off: %v vs. %v", name, av, bv)
				return false
			}
		case rune:
			if av.(rune) != bv.(rune) {
				t.Errorf("rune field %v was wrong: %c != %c", name, av, bv)
				return false
			}
		case float64:
			if !near(av.(float64), bv.(float64)) {
				t.Errorf("Not close enough on field %v: %v vs. %v", name, av, bv)
				return false
			}
		default:
			if !reflect.DeepEqual(av, bv) {
				t.Errorf("%T field %v was wrong: %v != %v", av, name, av, bv)
				return false
			}
		}
	}

	return true
}

func TestRMCHandling(t *testing.T) {
	h := &rmcHandler{}
	for _, s := range strings.Split(ubloxSample, "\n") {
		parseMessage(s, h)
	}
	exp := RMC{
		Timestamp: time.Unix(1152634974, 0).UTC(),
		Status:    'A',
		Latitude:  37.383806166666666,
		Longitude: -121.9899755,
		Speed:     0.82,
		Angle:     188.36,
		Magvar:    0,
	}
	if !similar(t, h.rmc, exp) {
		t.Errorf("Expected more similarity between %#v and (wanted) %#v", h.rmc, exp)
	}
}

type vtgHandler struct {
	vtg VTG
}

func (r *vtgHandler) HandleVTG(vtg VTG) {
	r.vtg = vtg
}

func TestVTGHandling(t *testing.T) {
	h := &vtgHandler{}
	for _, s := range strings.Split(ubloxSample, "\n") {
		parseMessage(s, h)
	}
	exp := VTG{
		True:     188.36,
		Magnetic: 0,
		Knots:    0.82,
		KMH:      1.519,
	}
	if !similar(t, h.vtg, exp) {
		t.Errorf("Expected more similarity between %#v and (wanted) %#v", h.vtg, exp)
	}
}

type ggaHandler struct {
	gga GGA
}

func (g *ggaHandler) HandleGGA(gga GGA) {
	g.gga = gga
}

func TestFixQualityStringing(t *testing.T) {
	got := fmt.Sprint(FloatRealTimeKinematicFix)
	if got != "float rt kinematic" {
		t.Errorf("Incorrect value for FloatRealTimeKinematicFix: %v", got)
	}
}

func TestGGAHandling(t *testing.T) {
	h := &ggaHandler{}
	for _, s := range strings.Split(ubloxSample, "\n") {
		parseMessage(s, h)
	}
	exp := GGA{
		Taken:              time.Date(0, 1, 1, 16, 22, 54, 0, time.UTC),
		Latitude:           37.383806166666666,
		Longitude:          -121.9899755,
		Quality:            GPSFix,
		NumSats:            3,
		HorizontalDilution: 2.36,
		Altitude:           525.6,
		GeoidHeight:        -25.6,
	}
	if !similar(t, h.gga, exp) {
		t.Errorf("Expected more similarity between %#v and (wanted) %#v", h.gga, exp)
	}
}

type gsaHandler struct {
	gsa GSA
}

func (g *gsaHandler) HandleGSA(gsa GSA) {
	g.gsa = gsa
}

func TestGSAHandling(t *testing.T) {
	h := &gsaHandler{}
	for _, s := range strings.Split(ubloxSample, "\n") {
		parseMessage(s, h)
	}
	exp := GSA{
		Auto:     true,
		Fix:      Fix2D,
		SatsUsed: []int{25, 1, 22},
		PDOP:     2.56,
		HDOP:     2.36,
		VDOP:     1,
	}
	if !similar(t, h.gsa, exp) {
		t.Errorf("Expected more similarity between %#v and (wanted) %#v", h.gsa, exp)
	}
}

type gllHandler struct {
	gll GLL
}

func (g *gllHandler) HandleGLL(gll GLL) {
	g.gll = gll
}

func TestGLLHandling(t *testing.T) {
	h := &gllHandler{}
	for _, s := range strings.Split(ubloxSample, "\n") {
		parseMessage(s, h)
	}
	exp := GLL{
		Latitude:  37.383806166666666,
		Longitude: -121.9899755,
		Active:    true,
		Taken:     time.Date(0, 1, 1, 16, 22, 54, 0, time.UTC),
	}
	if !similar(t, h.gll, exp) {
		t.Errorf("Expected more similarity between %#v and (wanted) %#v", h.gll, exp)
	}
}

type zdaHandler struct {
	zda ZDA
}

func (g *zdaHandler) HandleZDA(zda ZDA) {
	g.zda = zda
}

// $GPZDA,162254.00,11,07,2006,00,00*63
func TestZDAHandling(t *testing.T) {
	h := &zdaHandler{}
	for _, s := range strings.Split(ubloxSample, "\n") {
		parseMessage(s, h)
	}
	exp := ZDA{time.Date(2006, 7, 11, 16, 22, 54, 0, time.UTC)}
	if !similar(t, h.zda, exp) {
		t.Errorf("Expected more similarity between %#v and (wanted) %#v", h.zda, exp)
	}
}

func TestZDAZones(t *testing.T) {
	tests := map[string]time.Time{
		"$GPZDA,162254.00,11,07,2006,00,00*63": time.Date(2006, 7, 11, 16, 22, 54, 0, time.UTC),
		"$GPZDA,050306,29,10,2003,,*43":        time.Date(2003, 10, 29, 5, 3, 6, 0, time.UTC),
		"$GPZDA,110003.00,27,03,2006,-5,00*7f": time.Date(2006, 3, 27, 11, 0, 3, 0, time.FixedZone("GPS", -18000)),
	}

	for in, exp := range tests {
		h := &zdaHandler{}
		parseMessage(in, h)
		if !similar(t, h.zda, ZDA{exp}) {
			t.Errorf("Expected more similarity between %#v and (wanted) %#v", h.zda, exp)
		}

	}
}

type gsvHandler struct {
	gsv GSV
}

func (g *gsvHandler) HandleGSV(gsv GSV) {
	g.gsv = gsv
}

func TestGSVAccumulation(t *testing.T) {
	in := []GSV{
		// Send a few out of order
		{TotalSentences: 4, SentenceNum: 2, InView: 14, SatInfo: []GSVSatInfo{
			{18, 16, 79, 0},
			{11, 19, 312, 0},
			{14, 80, 41, 0},
			{21, 4, 135, 25},
		}},
		{TotalSentences: 4, SentenceNum: 1, InView: 14, SatInfo: []GSVSatInfo{
			{25, 15, 175, 30},
			{14, 80, 41, 0},
			{19, 38, 259, 14},
			{1, 52, 233, 18},
		}},
		{TotalSentences: 4, SentenceNum: 3, InView: 14, SatInfo: []GSVSatInfo{
			{15, 27, 134, 18},
			{3, 25, 222, 0},
			{22, 51, 57, 16},
			{9, 7, 36, 0},
		}},

		// Now the real ones
		{TotalSentences: 4, SentenceNum: 1, InView: 14, SatInfo: []GSVSatInfo{
			{25, 15, 175, 30},
			{14, 80, 41, 0},
			{19, 38, 259, 14},
			{1, 52, 233, 18},
		}},
		{TotalSentences: 4, SentenceNum: 2, InView: 14, SatInfo: []GSVSatInfo{
			{18, 16, 79, 0},
			{11, 19, 312, 0},
			{14, 80, 41, 0},
			{21, 4, 135, 25},
		}},
		{TotalSentences: 4, SentenceNum: 3, InView: 14, SatInfo: []GSVSatInfo{
			{15, 27, 134, 18},
			{3, 25, 222, 0},
			{22, 51, 57, 16},
			{9, 7, 36, 0},
		}},
		{TotalSentences: 4, SentenceNum: 4, InView: 14, SatInfo: []GSVSatInfo{
			{7, 1, 181, 0},
			{15, 25, 135, 0},
		}},
	}
	exp := GSVAccumulator{
		InView: 14,
		Parts:  4,
		prev:   4,
		SatInfo: []GSVSatInfo{
			{25, 15, 175, 30},
			{14, 80, 41, 0},
			{19, 38, 259, 14},
			{1, 52, 233, 18},
			{18, 16, 79, 0},
			{11, 19, 312, 0},
			{14, 80, 41, 0},
			{21, 4, 135, 25},
			{15, 27, 134, 18},
			{3, 25, 222, 0},
			{22, 51, 57, 16},
			{9, 7, 36, 0},
			{7, 1, 181, 0},
			{15, 25, 135, 0},
		},
	}

	a := GSVAccumulator{}
	for _, g := range in {
		a.Add(g)
	}

	if !similar(t, a, exp) {
		t.Errorf("Expected more similarity between %#v and (wanted) %#v", a, exp)
	}
}

// $GPGSV,4,1,14, 25,15,175,30, 14,80,041,,  19,38,259,14,  01,52,223,18   *76
// $GPGSV,4,2,14, 18,16,079,,   11,19,312,,  14,80,041,,    21,04,135,25   *7D
// $GPGSV,4,3,14, 15,27,134,18, 03,25,222,,  22,51,057,16,  09,07,036,     *79
// $GPGSV,4,4,14, 07,01,181,,   15,25,135,                                 *76
func TestGSVHandling(t *testing.T) {
	h := &gsvHandler{}
	for _, s := range strings.Split(ubloxSample, "\n") {
		parseMessage(s, h)
	}

	exp := GSV{
		InView:         14,
		SentenceNum:    4,
		TotalSentences: 4,
		SatInfo: []GSVSatInfo{
			{7, 1, 181, 0},
			{15, 25, 135, 0},
		},
	}
	if !similar(t, h.gsv, exp) {
		t.Errorf("Expected more similarity between %#v and (wanted) %#v", h.gsv, exp)
	}
}

type aamHandler struct {
	aam AAM
}

func (g *aamHandler) HandleAAM(aam AAM) {
	g.aam = aam
}

// $GPAAM,A,A,0.10,N,WPTNME*32
func TestAAMHandling(t *testing.T) {
	h := &aamHandler{}
	for _, s := range strings.Split(ubloxSample, "\n") {
		parseMessage(s, h)
	}
	exp := AAM{
		Arrival:       true,
		Perpendicular: true,
		Radius:        0.1,
	}
	if !similar(t, h.aam, exp) {
		t.Errorf("Expected more similarity between %#v and (wanted) %#v", h.aam, exp)
	}
}

type gstHandler struct {
	gst GST
}

func (g *gstHandler) HandleGST(gst GST) {
	g.gst = gst
}

// $GPGST,024603.00,3.2,6.6,4.7,47.3,5.8,5.6,22.0*58
func TestGSTHandling(t *testing.T) {
	h := &gstHandler{}
	for _, s := range strings.Split(ubloxSample, "\n") {
		parseMessage(s, h)
	}
	exp := GST{
		Timestamp:             time.Date(0, 1, 1, 02, 46, 03, 0, time.UTC),
		Deviation:             3.2,
		MajorDeviceation:      6.6,
		MinorDeviation:        4.7,
		MajorOrientation:      47.3,
		MinorOrientation:      5.8,
		LatitudeErrDeviation:  5.6,
		LongitudeErrDeviation: 22,
	}
	if !similar(t, h.gst, exp) {
		t.Errorf("Expected more similarity between %#v and (wanted) %#v", h.gst, exp)
	}
}
