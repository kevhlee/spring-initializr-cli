const qualifiers = ["M", "RC", "BUILD-SNAPSHOT", "RELEASE"];
const rangePattern = /(\(|\[)(.*),(.*)(\)|\])/;
const versionPattern = /^(\d+)\.(\d+|x)\.(\d+|x)(?:([.|-])([^0-9]+)(\d+)?)?$/;

export function withinRange(range: string, version: string): boolean {
  const match = range.match(rangePattern);
  if (!match) {
    return true;
  }

  const lowerInclusive = match[1] === "[";
  const higherInclusive = match[4] === "]";

  const lowerVersion = match[2];
  const higherVersion = match[3];

  if (lowerInclusive && higherInclusive) {
    return (
      compareVersion(lowerVersion, version) <= 0 &&
      compareVersion(higherVersion, version) >= 0
    );
  } else if (lowerInclusive) {
    return (
      compareVersion(lowerVersion, version) <= 0 &&
      compareVersion(higherVersion, version) > 0
    );
  } else if (higherInclusive) {
    return (
      compareVersion(lowerVersion, version) < 0 &&
      compareVersion(higherVersion, version) >= 0
    );
  } else {
    return (
      compareVersion(lowerVersion, version) < 0 &&
      compareVersion(higherVersion, version) > 0
    );
  }
}

export function compareVersion(a: string, b: string): number {
  const versionA = a.split(".");
  const versionB = b.split(".");

  for (let i = 0; i < 3; i++) {
    const va = parseInt(versionA[i]);
    const vb = parseInt(versionB[i]);
    const result = va - vb;

    if (result !== 0) {
      return result;
    }
  }

  return compareQualifier(parseQualifier(a), parseQualifier(b));
}

export function compareQualifier(a: string, b: string): number {
  let indexA = qualifiers.length - 1;
  let indexB = qualifiers.length - 1;

  for (let i = 0; i < qualifiers.length; i++) {
    if (a === qualifiers[i]) {
      indexA = i;
    }
    if (b === qualifiers[i]) {
      indexB = i;
    }
  }

  return indexA - indexB;
}

export function parseQualifier(version: string): string {
  const match = version.match(versionPattern);
  if (!match) {
    return "RELEASE";
  }

  const qualifier = match[5];
  if (!qualifier) {
    return "RELEASE";
  }

  for (const q of qualifiers) {
    if (q === qualifier) {
      return qualifier;
    }
  }

  return "RELEASE";
}
