import AdmZip from "adm-zip";
import fetch from "node-fetch";

const DEFAULT_ACCEPT_HEADER = "application/vnd.initializr.v2.2+json";

export type InitializrMetadata = {
  artifactId: InitializrMetadataText;
  bootVersion: InitializrMetadataSelect;
  dependencies: InitializrMetadataMultiSelect;
  description: InitializrMetadataText;
  groupId: InitializrMetadataText;
  javaVersion: InitializrMetadataSelect;
  language: InitializrMetadataSelect;
  name: InitializrMetadataText;
  packageName: InitializrMetadataText;
  packaging: InitializrMetadataSelect;
  type: InitializrMetadataSelect;
  version: InitializrMetadataText;
};

export type InitializrMetadataType =
  | "action"
  | "hierarchical-multi-select"
  | "single-select"
  | "text";

export type InitializrMetadataMultiSelect = {
  values: InitializrMetadataSelectGroup[];
};

export type InitializrMetadataSelect = {
  type: InitializrMetadataType;
  default: string;
  values: InitializrMetadataSelectValue[];
};

export type InitializrMetadataSelectGroup = {
  name: string;
  values: InitializrMetadataSelectValue[];
};

export type InitializrMetadataSelectValue = {
  id: string;
  name: string;
  description?: string;
  versionRange?: string;
};

export type InitializrMetadataText = {
  type: InitializrMetadataType;
  default: string;
};

export type InitializrMetadataElement =
  | InitializrMetadataMultiSelect
  | InitializrMetadataSelect
  | InitializrMetadataText;

export async function fetchInitializrMetadata(
  url: string
): Promise<InitializrMetadata> {
  const response = await fetch(url, {
    method: "GET",
    headers: {
      Accept: DEFAULT_ACCEPT_HEADER,
    },
  });

  return (await response.json()) as InitializrMetadata;
}

export type InitializrParameters = { [id: string]: string | string[] };

export async function generateInitializrProject(
  url: string | URL,
  parameters: InitializrParameters,
  zipped: boolean = false
) {
  const generateURL = new URL("starter.zip", url);

  for (const id of Object.keys(parameters)) {
    let value: string | string[] = parameters[id];

    if (Array.isArray(value)) {
      value = value.join(",");
    }

    if (value.length > 0) {
      generateURL.searchParams.set(id, value);
    }
  }

  const response = await fetch(generateURL, {
    method: "GET",
    headers: {
      Accept: DEFAULT_ACCEPT_HEADER,
    },
  });

  const zip: AdmZip = new AdmZip(await response.buffer());

  if (zipped) {
    zip.writeZip(parameters.artifactId + ".zip");
  } else {
    zip.extractAllTo(parameters.artifactId as string, true, true);
  }
}
