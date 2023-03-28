package egts

import (
	"errors"
)

const (
	EgtsPcOk             = uint8(0)   // код сообщения, что пакет успешно обработано
	EgtsPcInProgress     = uint8(1)   // код сообщения, что пакет в процессе обработки
	EgtsPcUnsProtocol    = uint8(128) // неподдерживаемый протокол
	EgtsPcDecryptError   = uint8(129) // ошибка декодирования
	EgtsPcProcDenied     = uint8(130) // обработка запрещена
	EgtsPcIncHeaderform  = uint8(131) // неверный формат заголовка
	EgtsPcIncDataform    = uint8(132) // неверный формат данных
	EgtsPcUnsType        = uint8(133) // EgtsPcUnsType неподдерживаемый тип
	EgtsPcNotenParams    = uint8(134) // неверное количество параметров
	EgtsPcDblProc        = uint8(135) // попытка повторной обработки
	EgtsPcProcSrcDenied  = uint8(136) // обработка данных от источника запрещена
	EgtsPcHeaderCrcError = uint8(137) // ошибка контрольной суммы заголовка
	EgtsPcDatacrcError   = uint8(138) // ошибка контрольной суммы данных
	EgtsPcInvdatalen     = uint8(139) // некорректная длина данных
	EgtsPcRouteNfound    = uint8(140) // маршрут не найден
	EgtsPcRouteClosed    = uint8(141) // маршрут закрыт
	EgtsPcRouteDenied    = uint8(142) // маршрутизация запрещена
	EgtsPcInvaddr        = uint8(143) // неверный адрес
	EgtsPcTtlexpired     = uint8(144) // превышено количество ретрансляции данных
	EgtsPcNoAck          = uint8(145) // нет подтверждения
	EgtsPcObjNfound      = uint8(146) // объект не найден
	EgtsPcEvntNfound     = uint8(147) // событие не найдено
	EgtsPcSrvcNfound     = uint8(148) // сервис не найден
	EgtsPcSrvcDenied     = uint8(149) // сервис запрещён
	EgtsPcSrvcUnkn       = uint8(150) // неизвестный тип сервиса
	EgtsPcAuthPenied     = uint8(151) // авторизация запрещена
	EgtsPcAlreadyExists  = uint8(152) // объект уже существует
	EgtsPcIDNfound       = uint8(153) // идентификатор не найден
	EgtsPcIncDatetime    = uint8(154) // неправильная дата и время
	EgtsPcIoError        = uint8(155) // ошибка ввода/вывода
	EgtsPcNoResAvail     = uint8(156) // недостаточно ресурсов
	EgtsPcModuleFault    = uint8(157) // внутренний сбой модуля
	EgtsPcModulePwrFlt   = uint8(158) // сбой в работе цепи питания модуля
	EgtsPcModuleProcFlt  = uint8(159) // сбой в работе микроконтроллера модуля
	EgtsPcModuleSwFlt    = uint8(160) // сбой в работе программы модуля
	EgtsPcModuleFwFlt    = uint8(161) // сбой в работе внутреннего ПО модуля
	EgtsPcModuleIoFlt    = uint8(162) // сбой в работе блока ввода/вывода модуля
	EgtsPcModuleMemFlt   = uint8(163) // сбой в работе внутренней памяти модуля
	EgtsPcTestFailed     = uint8(164) // тест не пройден
)

var (
	// ErrSecretKey represents the error of secret key is nil.
	ErrSecretKey = errors.New("package is encrypted but secret key is nil")
)
