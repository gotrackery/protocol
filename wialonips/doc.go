// Packet wialonips implements the Wialon IPS communication protocol that was developed by Gurtam for use in personal
// and automotive GPS and GLONASS trackers which transfer data to a satellite monitoring server using the TCP or the
// UDP protocol.
//
// The TCP connection must be maintained throughout the entire data transfer process.
// If the device disconnects immediately after sending the message,
// the server does not have time to send a response to the device, and traffic consumption increases.
// While using one TCP connection, you should transfer data from one device. Otherwise,
// the system registers only the data of the device whose ID is the first in the incoming data list.
// To save traffic, you can use the UDP protocol. However, it does not guarantee that the messages will be delivered.
package wialonips
