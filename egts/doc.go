// Package egts implements EGTS (Era Glonass Telematics Standard) - telematic data transmission standard developed by
// the Government of the Russian Federation and approved by Order #285 of the Ministry of Transport of July 31, 2012
//
// http://docs.cntd.ru/document/1200095098
// EGTS_PACKET {
// _EGTS_PACKET_HEADER
// __SFRD - data structure, depending on the type of package EGTS_PACKET_HEADER.PT
// __SFRCS - checksum SFRD CRC-16
// }
// The total length of the Transport Layer Protocol packet does not exceed 65535 bytes
// https://github.com/LdDl/go-egts/blob/master/docs_rus/egts.txt
package egts
