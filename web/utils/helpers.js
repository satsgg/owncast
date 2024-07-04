import UAParser from 'ua-parser-js';

export function pluralize(string, count) {
  if (count === 1) {
    return string;
  }
  return `${string}s`;
}

export function getDiffInDaysFromNow(timestamp) {
  const time = typeof timestamp === 'string' ? new Date(timestamp) : timestamp;
  return (new Date() - time) / (24 * 3600 * 1000);
}

// Take a nested object of state metadata and merge it into
// a single flattened node.
export function mergeMeta(meta) {
  return Object.keys(meta).reduce((acc, key) => {
    const value = meta[key];
    Object.assign(acc, value);

    return acc;
  }, {});
}

export const isMobileSafariIos = () => {
  try {
    const ua = navigator.userAgent;
    const uaParser = new UAParser(ua);
    const browser = uaParser.getBrowser();
    const device = uaParser.getDevice();

    if (device.vendor !== 'Apple') {
      return false;
    }

    if (device.type !== 'mobile' && device.type !== 'tablet') {
      return false;
    }

    return browser.name === 'Mobile Safari' || browser.name === 'Safari';
  } catch (e) {
    return false;
  }
};

export const isMobileSafariHomeScreenApp = () => {
  if (!isMobileSafariIos()) {
    return false;
  }

  return 'standalone' in window.navigator && window.navigator.standalone;
};

export function base64ToHex(base64) {
  // Decode base64 to binary string
  const binaryString = atob(base64);

  // Convert binary string to hex
  let hex = '';
  for (let i = 0; i < binaryString.length; i++) {
    const hexChar = binaryString.charCodeAt(i).toString(16);
    hex += hexChar.padStart(2, '0'); // Ensure each hex value is two digits
  }

  return hex;
}
