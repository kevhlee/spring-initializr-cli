import { Command, program } from "commander";
import figlet from "figlet";
import prompts, {
  Answers,
  Choice,
  InitialReturnValue,
  PrevCaller,
  PromptObject,
} from "prompts";
import {
  fetchInitializrMetadata,
  generateInitializrProject,
  InitializrMetadata,
  InitializrMetadataMultiSelect,
  InitializrMetadataSelect,
  InitializrMetadataText,
  InitializrParameters,
} from "./initializr";
import { withinRange } from "./version";

export function runCLI(argv: string[]): Promise<Command> {
  return program
    .name("initializr")
    .description("A CLI for creating Spring projects locally")
    .version("v1.0.0")
    .addHelpText(
      "beforeAll",
      figlet.textSync("Initializr CLI", { font: "Small Slant" })
    )
    .argument(
      "[url]",
      "the URL of the Spring Initializr instance",
      "https://start.spring.io"
    )
    .option("--zipped", "keep generated project zipped", false)
    .action(async (url: string, options: { zipped: boolean }) => {
      const metadata = await fetchInitializrMetadata(url);
      const parameters = await initiatePrompts(metadata);
      await generateInitializrProject(url, parameters, options.zipped);

      console.log("\nSpring project successfully generated. 🌱");
    })
    .parseAsync(argv);
}

async function initiatePrompts(
  metadata: InitializrMetadata
): Promise<InitializrParameters> {
  const questions: PromptObject<string>[] = [];

  questions.push(createSelectPrompt("type", "Project: ", metadata.type));

  questions.push(
    createSelectPrompt("language", "Language: ", metadata.language)
  );

  questions.push(
    createSelectPrompt("bootVersion", "Spring Boot: ", metadata.bootVersion)
  );

  questions.push(createTextPrompt("groupId", "Group: ", metadata.groupId));

  questions.push(
    createTextPrompt("artifactId", "Artifact: ", metadata.artifactId)
  );

  questions.push(createTextPrompt("version", "Version: ", metadata.version));

  questions.push(
    createTextPrompt(
      "name",
      "Name: ",
      metadata.name,
      (_, answers: Answers<string>) => {
        return answers.artifactId || metadata.name.default;
      }
    )
  );

  questions.push(
    createTextPrompt("description", "Description: ", metadata.description)
  );

  questions.push(
    createTextPrompt(
      "packageName",
      "Package name: ",
      metadata.packageName,
      (_, answers: Answers<string>) => {
        const groupId: string = answers.groupId;
        const artifactId: string = answers.artifactId;

        if (!(groupId && artifactId)) {
          return metadata.artifactId.default;
        }

        return groupId + "." + artifactId;
      }
    )
  );

  questions.push(
    createSelectPrompt("packaging", "Packaging: ", metadata.packaging)
  );

  questions.push(
    createSelectPrompt("javaVersion", "Java: ", metadata.javaVersion)
  );

  questions.push(
    createMultiSelectPrompt(
      "dependencies",
      "Dependencies: ",
      metadata.dependencies
    )
  );

  const parameters = await prompts(questions, {
    onCancel: () => {
      throw new Error("Cancelled");
    },
  });

  const confirmResponse = await prompts({
    type: "confirm",
    name: "confirm",
    message: "Generate project?",
  });

  if (!confirmResponse.confirm) {
    throw new Error("Cancelled");
  }

  return parameters;
}

function createMultiSelectPrompt(
  name: string,
  message: string,
  multiselect: InitializrMetadataMultiSelect
): PromptObject<string> {
  return {
    type: "autocompleteMultiselect",
    name: name,
    message: message,
    choices: (_, answers: Answers<string>) => {
      const choices: Choice[] = [];
      const bootVersion: string = answers.bootVersion || "";

      for (const group of multiselect.values) {
        for (const value of group.values) {
          choices.push({
            title: value.name,
            value: value.id,
            description: value.description,
            disabled: !withinRange(value.versionRange || "", bootVersion),
          });
        }
      }

      choices.sort((a, b) =>
        a.title === b.title ? 0 : a.title < b.title ? -1 : 1
      );

      return choices;
    },
  };
}

function createSelectPrompt(
  name: string,
  message: string,
  select: InitializrMetadataSelect
): PromptObject<string> {
  const choices = select.values.map((value) => {
    return {
      title: value.name,
      value: value.id,
      description: value.description,
    };
  });

  const initial = choices.findIndex((choice) => {
    return choice.value === select.default;
  });

  return {
    type: "select",
    name: name,
    message: message,
    choices: choices,
    initial: initial,
  };
}

function createTextPrompt(
  name: string,
  message: string,
  text: InitializrMetadataText,
  initial?: PrevCaller<string, InitialReturnValue>
): PromptObject<string> {
  return {
    type: "text",
    name: name,
    message: message,
    initial: initial || text.default,
  };
}
