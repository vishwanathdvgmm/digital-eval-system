import React, { useState } from "react";
import { uploadScript } from "../../api/examiner";
import { ScriptUploadResponse } from "../../types/examiner";
import Card from "../../components/Card";

const UploadScripts: React.FC = () => {
    const [file, setFile] = useState<File | null>(null);
    const [uploading, setUploading] = useState(false);
    const [result, setResult] = useState<ScriptUploadResponse | null>(null);
    const [error, setError] = useState<string | null>(null);

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files && e.target.files.length > 0) {
            setFile(e.target.files[0]);
            setError(null);
            setResult(null);
        }
    };

    const handleUpload = async () => {
        if (!file) {
            setError("Please select a file first.");
            return;
        }

        setUploading(true);
        setError(null);
        setResult(null);

        try {
            const data = await uploadScript(file);
            setResult(data);
        } catch (err: any) {
            console.error("Upload failed", err);
            setError(err.response?.data || "Upload failed. Please try again.");
        } finally {
            setUploading(false);
        }
    };

    return (
        <div className="space-y-6">
            <h1 className="text-2xl font-bold text-gray-800">Upload Answer Script</h1>

            <Card>
                <div className="space-y-4">
                    <div className="border-2 border-dashed border-gray-300 rounded-lg p-8 text-center hover:border-blue-500 transition-colors">
                        <input
                            type="file"
                            accept=".pdf"
                            onChange={handleFileChange}
                            className="hidden"
                            id="file-upload"
                        />
                        <label
                            htmlFor="file-upload"
                            className="cursor-pointer flex flex-col items-center justify-center"
                        >
                            <svg
                                className="w-12 h-12 text-gray-400 mb-3"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                    strokeWidth={2}
                                    d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
                                />
                            </svg>
                            <span className="text-gray-600 font-medium">
                                {file ? file.name : "Click to select PDF file"}
                            </span>
                            <span className="text-xs text-gray-400 mt-1">
                                Only PDF files are allowed
                            </span>
                        </label>
                    </div>

                    {error && (
                        <div className="p-3 bg-red-100 text-red-700 rounded-md text-sm">
                            {error}
                        </div>
                    )}

                    <button
                        onClick={handleUpload}
                        disabled={!file || uploading}
                        className={`w-full py-2 px-4 rounded-md text-white font-medium transition-colors ${!file || uploading
                                ? "bg-gray-400 cursor-not-allowed"
                                : "bg-blue-600 hover:bg-blue-700"
                            }`}
                    >
                        {uploading ? "Uploading & Extracting..." : "Upload Script"}
                    </button>
                </div>
            </Card>

            {result && (
                <Card className="border-l-4 border-l-green-500 bg-green-50">
                    <h2 className="text-lg font-semibold text-green-800 mb-2">
                        Upload Successful!
                    </h2>
                    <div className="space-y-2 text-sm text-green-700">
                        <p>
                            <span className="font-semibold">Script ID:</span> {result.script_id}
                        </p>
                        <p>
                            <span className="font-semibold">Block Hash:</span> {result.block_hash}
                        </p>
                        <p>
                            <span className="font-semibold">IPFS CID:</span> {result.pdf_cid}
                        </p>
                        <div className="mt-2 pt-2 border-t border-green-200">
                            <p className="font-semibold mb-1">Extracted Metadata:</p>
                            <ul className="list-disc list-inside pl-2">
                                {Object.entries(result.metadata).map(([key, value]) => (
                                    <li key={key}>
                                        {key}: {value}
                                    </li>
                                ))}
                            </ul>
                        </div>
                    </div>
                </Card>
            )}
        </div>
    );
};

export default UploadScripts;
