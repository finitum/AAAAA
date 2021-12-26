import humanizeDuration from "@/utils/humanizeDuration";

export interface Package {
  // Name is the name of the package (unique)
  Name: string;
  // RepoURL is the git repository where the PKGBUILD can be found
  RepoURL: string;
  // RepoBranch is which branch is used for updates
  RepoBranch: string;
  // KeepLastN determines how many old versions of packages are kept
  KeepLastN: number;
  // UpdateFrequency determines how often the package should be updated
  UpdateFrequency: number;
}

export function NewPackage(): Package {
  return {
    KeepLastN: 2,
    Name: "",
    RepoBranch: "",
    RepoURL: "",
    UpdateFrequency: 0
  };
}

export function frequencyToDuration(freqns: number): string {
  return humanizeDuration(freqns / 1000 / 1000);
}

export interface User {
  Username: string;
  Password: string;
}

export interface Job {
  PackageName: string;
  Status: BuildStatus;
  Logs: LogLine[];
  Uuid: string;
  Time: number;
}

export enum BuildStatus {
  Pending,
  PullingRepo,
  Running,
  Uploading,
  Done,

  Errored
}

export function BuildStatusToString(status: BuildStatus): string {
  switch (status) {
    case BuildStatus.Pending:
      return "Queued";
    case BuildStatus.PullingRepo:
      return "Pulling repository";
    case BuildStatus.Running:
      return "Running";
    case BuildStatus.Uploading:
      return "Uploading";
    case BuildStatus.Done:
      return "Done";
    case BuildStatus.Errored:
      return "Errored";

    default:
      return "Unknown";
  }
}

export enum LogLevel {
  PanicLevel,
  FatalLevel,
  ErrorLevel,
  WarnLevel,
  InfoLevel,
  DebugLevel,
  TraceLevel
}

export interface LogLine {
  Time: number;
  Level: LogLevel;
  message: string;
}
