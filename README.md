# vpn
A TUI wg-quick interface.

# usage
Running without any arguments launches the TUI interface. Pressing j and k will
navigate the list. Enter will select a wireguard config. Frequently used configs
will be at the top of the list.

A simple cli mode also exists for quickly connecting to your most used config or
stopping the current connection.
```
vpn [ on | up | off | down ]
```

Both modes read a list of configured servers in /etc/wireguard using the format
described by wg-quick.

# authors
"Maintained" by Dakota Walsh <kota at nilsu.org>.
Up-to-date sources can be found at https://git.sr.ht/~kota/vpn/

## license
GNU GPL version 3 or later, see LICENSE.
