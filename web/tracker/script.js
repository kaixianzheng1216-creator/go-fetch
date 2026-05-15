(function () {
  var currentScript = document.currentScript;
  if (!currentScript) return;

  var website = currentScript.getAttribute("data-website-id");
  var distinctId = currentScript.getAttribute("data-distinct-id") || "";
  var endpoint = currentScript.getAttribute("data-host-url") || "";
  endpoint = endpoint.replace(/\/$/, "") + "/api/collect";

  var screenSize = window.screen
    ? window.screen.width + "x" + window.screen.height
    : "";
  var currentUrl = location.href;
  var currentRef =
    document.referrer &&
    new URL(document.referrer, location.href).hostname === location.hostname
      ? ""
      : document.referrer;

  function payload(name, data) {
    return {
      website: website,
      url: currentUrl,
      referrer: currentRef,
      title: document.title,
      screen: screenSize,
      language: navigator.language,
      distinctId: distinctId || undefined,
      name: name || "",
      data: data || undefined,
    };
  }

  function send(name, data) {
    if (!website) return;
    var body = JSON.stringify({ type: "event", payload: payload(name, data) });

    try {
      navigator.sendBeacon && !data
        ? navigator.sendBeacon(
            endpoint,
            new Blob([body], { type: "application/json" }),
          )
        : fetch(endpoint, {
            method: "POST",
            keepalive: true,
            headers: { "Content-Type": "application/json" },
            body: body,
          });
    } catch (_) {}
  }

  function track(name, data) {
    if (typeof name === "string") return send(name, data);
    return send("", undefined);
  }

  function routeChanged(url) {
    if (!url) return;
    var next = new URL(url, location.href).href;
    if (next === currentUrl) return;
    currentRef = currentUrl;
    currentUrl = next;
    setTimeout(function () {
      track();
    }, 100);
  }

  function hook(method) {
    var original = history[method];
    history[method] = function () {
      var result = original.apply(this, arguments);
      routeChanged(arguments[2]);
      return result;
    };
  }

  window.goFetch = window.goFetch || {};
  window.goFetch.track = track;
  hook("pushState");
  hook("replaceState");
  window.addEventListener("popstate", function () {
    routeChanged(location.href);
  });

  if (document.readyState === "complete") {
    track();
  } else {
    window.addEventListener("load", function () {
      track();
    });
  }
})();
