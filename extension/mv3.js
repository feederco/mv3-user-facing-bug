const CONFIG = {
  vapidPublicKey: "BAzsP8FJ4nf_fPgTTv8Agj5z6WbIJMFWr7AezO3_b_zfLuCFhrzE8O1GLRvfXKQ7B4JKkElxlLBKjEszW7NuYQc",
};

// UNCOMMENT THIS TO FIX THE ISSUE
/*registration.pushManager
  .subscribe({
    userVisibleOnly: false,
    applicationServerKey: urlB64ToUint8Array(CONFIG.vapidPublicKey),
  })
  .then((subscription) => {})
  .catch((e) => {
    console.error("[Failed to subscribe to push service", e);
  });
*/

chrome.runtime.onStartup.addListener(async (e) => {
  console.log("chrome.runtime.onStartup", e);
  // Throws an error "service worker is not registered"
  // registerWebPush();
});

chrome.runtime.onInstalled.addListener((e) => {
  console.log("chrome.runtime.onInstalled", e);
  // Throws an error "service worker is not registered"
  // registerWebPush();
});

addEventListener("push", (e) => {
  console.log("PUSH CALLED!", e.data.text(), e);
});

addEventListener("activate", (e) => {
  console.log("addEventListener(activate) called");
  installWebPush();
});

async function installWebPush() {
  console.log("installWebPush");
  const registration = await registerWebPush();
  if (registration) {
    console.log("syncWebPushToken with", registration);
    syncWebPushToken(registration);
  }
}

async function registerWebPush() {
  try {
    const subscribe = await registration.pushManager.subscribe({
      userVisibleOnly: false,
      applicationServerKey: urlB64ToUint8Array(CONFIG.vapidPublicKey),
    });

    return JSON.stringify(subscribe);
  } catch (e) {
    console.error("Failed to subscribe to push service", e);
    return false;
  }
}
async function syncWebPushToken(pushData) {
  const res = await fetch("http://localhost:1337/token", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(pushData),
  });

  const data = await res.json();
  console.log("syncWebPushToken", data);
}

function urlB64ToUint8Array(base64String) {
  const padding = "=".repeat((4 - (base64String.length % 4)) % 4);
  const base64 = (base64String + padding).replace(/-/g, "+").replace(/_/g, "/");

  const rawData = atob(base64);
  const outputArray = new Uint8Array(rawData.length);

  for (let i = 0; i < rawData.length; ++i) {
    outputArray[i] = rawData.charCodeAt(i);
  }
  return outputArray;
}
