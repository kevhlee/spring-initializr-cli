import { runCLI } from "./cli";

runCLI(process.argv).catch((err: string | Error) => {
  if (err instanceof Error) {
    console.error(err.message);
  } else {
    console.error(err);
  }
  process.exit(1);
});
