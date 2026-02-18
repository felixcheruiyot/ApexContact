# About ApexContact

ApexContact is a gaming live streaming website that works with game promoters in boxing and car racing to bring live game events to its users. ApexContact is designed to be highly
scalable game streaming service with a paywall. Users pay to view events. ApexContact with partners (promoters and other stakeholders) schedule events and when it goes live, ApexContact
delivers the best event streaming experience.

Using Golang at the back (because of its concurrency nature) and VueJS at the front, ApexContact delivers the best streaming experience with zero latency and buffering. Customers pay through
partners like IntaSend (for M-Pesa STK Push) for streaming to start. To limit piracy, ApexContact conducts fraud checks by locking IPs, detecting VPN use, device fingerprinting, etc. Each subscribe event generates a unique
link/token that is tracked and can be opened only on one device and not multiple. This ensure no sharing of content/pirating and only paid users enjoy the stream. 
