package egts

// SubRecord Types.
const (
	SrRecordResponseType     byte = 0  // SrRecordResponseType is subrecord code of SR_RECORD_RESPONSE.
	SrTermIdentityType       byte = 1  // SrTermIdentityType is subrecord code of SR_TERM_IDENTITY.
	SrModuleDataType         byte = 2  // SrModuleDataType is subrecord code of SR_MODULE_DATA.
	SrDispatcherIdentityType byte = 5  // SrDispatcherIdentityType is subrecord code of SR_DISPATCHER_IDENTITY.
	SrAuthInfoType           byte = 7  // SrAuthInfoType is subrecord code of SR_AUTH_INFO.
	SrResultCodeType         byte = 9  // SrResultCodeType is subrecord code of SR_RESULT_CODE.
	SrEgtsPlusDataType       byte = 15 // SrEgtsPlusDataType is subrecord code of SR_PLUS_DATA.
	SrPosDataType            byte = 16 // SrPosDataType is subrecord code of SR_POS_DATA.
	SrExtPosDataType         byte = 17 // SrExtPosDataType is subrecord code of SR_EXT_POS_DATA.
	SrAdSensorsDataType      byte = 18 // SrAdSensorsDataType is subrecord code of SR_AD_SENSORS_DATA.
	SrCountersDataType       byte = 19 // SrCountersDataType is subrecord code of SR_COUNTERS_DATA.
	SrType20                 byte = 20 // SrType20 depending on the length may contain SR_STATE_DATA section (
	// if 5 bytes long) or SR_ACCEL_DATA.
	SrStateDataType          byte = 21 // SrStateDataType is subrecord code of SR_STATE_DATA.
	SrLoopinDataType         byte = 22 // SrLoopinDataType is subrecord code of SR_TERM_IDENTITY_TYPE ToDo check.
	SrAbsDigSensDataType     byte = 23 // SrAbsDigSensDataType is subrecord code of SR_ABS_DIG_SENS_DATA.
	SrAbsAnSensDataType      byte = 24 // SrAbsAnSensDataType is subrecord code of SR_ABS_AN_SENS_DATA.
	SrAbsCntrDataType        byte = 25 // SrAbsCntrDataType is subrecord code of SR_ABS_CNTR_DATA.
	SrAbsLoopinDataType      byte = 26 // SrAbsLoopinDataType is subrecord code of SR_ABS_LOOPIN_DATA.
	SrLiquidLevelSensorType  byte = 27 // SrLiquidLevelSensorType код is subrecord code of SR_LIQUID_LEVEL_SENSOR.
	SrPassengersCountersType byte = 28 // SrPassengersCountersType is subrecord code of SR_PASSENGERS_COUNTERS.
)

// Packet types.
const (
	// PtResponsePacket code of PT_RESPONSE type packet.
	PtResponsePacket byte = iota
	// PtAppdataPacket code of PT_APP_DATA type packet.
	PtAppdataPacket
	// PtSignedAppdataPacket code of PT_SIGNED_APPDATA type packet.
	PtSignedAppdataPacket
)

// Service types.
const (
	UndefinedService byte = iota
	// AuthService is service type of AUTH_SERVICE.
	AuthService
	// TeledataService is service type of TELEDATA_SERVICE.
	TeledataService
)
