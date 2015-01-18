package nmea

import "time"

// AAM represents a Waypoint Arrival Alarm message.
type AAM struct {
}

// A AAMHandler handles AAM messages from a stream.
type AAMHandler interface {
	HandleAAM(AAM)
}

// ALM represents a Almanac data message.
type ALM struct {
}

// A ALMHandler handles ALM messages from a stream.
type ALMHandler interface {
	HandleALM(ALM)
}

// APA represents a Auto Pilot A sentence message.
type APA struct {
}

// A APAHandler handles APA messages from a stream.
type APAHandler interface {
	HandleAPA(APA)
}

// APB represents a Auto Pilot B sentence message.
type APB struct {
}

// A APBHandler handles APB messages from a stream.
type APBHandler interface {
	HandleAPB(APB)
}

// BOD represents a Bearing Origin to Destination message.
type BOD struct {
}

// A BODHandler handles BOD messages from a stream.
type BODHandler interface {
	HandleBOD(BOD)
}

// BWC represents a Bearing using Great Circle route message.
type BWC struct {
}

// A BWCHandler handles BWC messages from a stream.
type BWCHandler interface {
	HandleBWC(BWC)
}

// DTM represents a Datum being used. message.
type DTM struct {
}

// A DTMHandler handles DTM messages from a stream.
type DTMHandler interface {
	HandleDTM(DTM)
}

// GGA represents a Fix information message.
type GGA struct {
}

// A GGAHandler handles GGA messages from a stream.
type GGAHandler interface {
	HandleGGA(GGA)
}

// GLL represents a Lat/Lon data message.
type GLL struct {
}

// A GLLHandler handles GLL messages from a stream.
type GLLHandler interface {
	HandleGLL(GLL)
}

// GRS represents a GPS Range Residuals message.
type GRS struct {
}

// A GRSHandler handles GRS messages from a stream.
type GRSHandler interface {
	HandleGRS(GRS)
}

// GSA represents a Overall Satellite data message.
type GSA struct {
}

// A GSAHandler handles GSA messages from a stream.
type GSAHandler interface {
	HandleGSA(GSA)
}

// GST represents a GPS Pseudorange Noise Statistics message.
type GST struct {
}

// A GSTHandler handles GST messages from a stream.
type GSTHandler interface {
	HandleGST(GST)
}

// GSV represents a Detailed Satellite data message.
type GSV struct {
}

// A GSVHandler handles GSV messages from a stream.
type GSVHandler interface {
	HandleGSV(GSV)
}

// MSK represents a send control for a beacon receiver message.
type MSK struct {
}

// A MSKHandler handles MSK messages from a stream.
type MSKHandler interface {
	HandleMSK(MSK)
}

// MSS represents a Beacon receiver status information message.
type MSS struct {
}

// A MSSHandler handles MSS messages from a stream.
type MSSHandler interface {
	HandleMSS(MSS)
}

// RMA represents a recommended Loran data message.
type RMA struct {
}

// A RMAHandler handles RMA messages from a stream.
type RMAHandler interface {
	HandleRMA(RMA)
}

// RMB represents a recommended navigation data for gps message.
type RMB struct {
}

// A RMBHandler handles RMB messages from a stream.
type RMBHandler interface {
	HandleRMB(RMB)
}

// RMC represents a recommended minimum data for gps message.
type RMC struct {
	Timestamp           time.Time
	Status              rune
	Latitude, Longitude float64
	Speed               float64
	Angle               float64
	Magvar              float64
}

// A RMCHandler handles RMC messages from a stream.
type RMCHandler interface {
	HandleRMC(RMC)
}

// RTE represents a route message message.
type RTE struct {
}

// A RTEHandler handles RTE messages from a stream.
type RTEHandler interface {
	HandleRTE(RTE)
}

// TRF represents a Transit Fix Data message.
type TRF struct {
}

// A TRFHandler handles TRF messages from a stream.
type TRFHandler interface {
	HandleTRF(TRF)
}

// STN represents a Multiple Data ID message.
type STN struct {
}

// A STNHandler handles STN messages from a stream.
type STNHandler interface {
	HandleSTN(STN)
}

// VBW represents a dual Ground / Water Spped message.
type VBW struct {
}

// A VBWHandler handles VBW messages from a stream.
type VBWHandler interface {
	HandleVBW(VBW)
}

// VTG represents a Vector track an Speed over the Ground message.
type VTG struct {
}

// A VTGHandler handles VTG messages from a stream.
type VTGHandler interface {
	HandleVTG(VTG)
}

// WCV represents a Waypoint closure velocity (Velocity Made Good) message.
type WCV struct {
}

// A WCVHandler handles WCV messages from a stream.
type WCVHandler interface {
	HandleWCV(WCV)
}

// WPL represents a Waypoint Location information message.
type WPL struct {
}

// A WPLHandler handles WPL messages from a stream.
type WPLHandler interface {
	HandleWPL(WPL)
}

// XTC represents a cross track error message.
type XTC struct {
}

// A XTCHandler handles XTC messages from a stream.
type XTCHandler interface {
	HandleXTC(XTC)
}

// XTE represents a measured cross track error message.
type XTE struct {
}

// A XTEHandler handles XTE messages from a stream.
type XTEHandler interface {
	HandleXTE(XTE)
}

// ZTG represents a Zulu (UTC) time and time to go (to destination) message.
type ZTG struct {
}

// A ZTGHandler handles ZTG messages from a stream.
type ZTGHandler interface {
	HandleZTG(ZTG)
}

// ZDA represents a Date and Time message.
type ZDA struct {
}

// A ZDAHandler handles ZDA messages from a stream.
type ZDAHandler interface {
	HandleZDA(ZDA)
}
