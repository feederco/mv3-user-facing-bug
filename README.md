# userFacing bug with Manifest V3 Chrome extensions

Reproduced as of 126.0.0

Steps:

1. Start server.go (this receives and stores the push tokens)
2. Install the extension
3. Close Chrome
4. Start Chrome and make sure to open a different user profile
5. Trigger a notification
6. See the "This site has been updated in the background" notification
7. Open the profile that has the extension installed
8. Trigger another notification
9. Don't see the background notification

Notifications are stored in `subscriptions.txt` in the server directory.

### Start Server

```bash
cd server
go run server.go
```

### Installing extension

1. Clone this repo
2. Open `chrome://extensions/` in Chrome
3. Enable Developer mode
4. Click "Load unpacked"
5. Select the `extension` directory

### Trigger a notification

```bash
cd server
go run server.go send
```

The script will also automatically clean up old subscriptions.

## Fix?

The fix appears to be to call `registration.pushManager.subscribe()` in the main runloop of the service worker. Calling that before the service worker is activated does throw an error, so it needs to be caught. But it does appear to silence the background notification.

Uncomment line 5 in `mv3.js` and try it again.
