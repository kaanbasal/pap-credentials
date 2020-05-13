# Pap Credential Logger

It requires a machine with 2 ethernet ports

- Connect cable coming from your ISP to one of the ethernet ports
- Connect cable coming from your router(WAN port) to the other ethernet port
- Run the application, and select the interfaces you just connected cables
- Wait for the result...

## How it works?

It reads packets from one network interface and relays all the packets without touching anything to the other, and vice versa.
If it captures any PAP packet containing username and password, it logs the information.

**Use at your own risk** 