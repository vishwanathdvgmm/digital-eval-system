export type MenuItem = {
  key: string;
  label: string;
  path: string;
  icon?: React.ReactNode;
};

export const roleMenu = (role: string | undefined): MenuItem[] => {
  switch (role) {
    case "authority":
      return [
        { key: "pending", label: "Pending Requests", path: "/authority/pending" },
        { key: "approve", label: "Approve Requests", path: "/authority/approve" },
        { key: "release", label: "Release Results", path: "/authority/release" },
      ];
    case "evaluator":
      return [
        { key: "assigned", label: "Assigned Scripts", path: "/evaluator/assigned" },
        { key: "request", label: "Request Evaluation", path: "/evaluator/request" },
      ];
    case "examiner":
      return [
        { key: "upload", label: "Upload Scripts", path: "/examiner/upload" },
        { key: "history", label: "Upload History", path: "/examiner/history" },
      ];
    case "student":
      return [
        { key: "view", label: "View Results", path: "/student/results" },
        { key: "download", label: "Download PDF", path: "/student/download" },
      ];
    default:
      return [{ key: "health", label: "Health", path: "/health" }];
  }
};