import React from "react";
import { ScriptRecord } from "../types/examiner";
import Card from "./Card";

interface Props {
    script: ScriptRecord;
}

const ScriptUploadItem: React.FC<Props> = ({ script }) => {
    return (
        <Card className="border-l-4 border-l-blue-500">
            <div className="flex justify-between items-start">
                <div>
                    <h3 className="text-lg font-semibold text-gray-800">
                        {script.course_id} - {script.usn}
                    </h3>
                    <p className="text-sm text-gray-500">
                        Semester: {script.semester} | Status:{" "}
                        <span
                            className={`font-medium ${script.status === "validated"
                                    ? "text-green-600"
                                    : script.status === "errored"
                                        ? "text-red-600"
                                        : "text-yellow-600"
                                }`}
                        >
                            {script.status.toUpperCase()}
                        </span>
                    </p>
                    <p className="text-xs text-gray-400 mt-1">ID: {script.script_id}</p>
                    <p className="text-xs text-gray-400">CID: {script.pdf_cid}</p>
                </div>
                <div className="text-right">
                    <p className="text-xs text-gray-400">
                        {new Date(script.created_at).toLocaleString()}
                    </p>
                </div>
            </div>
        </Card>
    );
};

export default ScriptUploadItem;
