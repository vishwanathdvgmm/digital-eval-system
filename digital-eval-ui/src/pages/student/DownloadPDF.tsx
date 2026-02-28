import React from "react";
import Card from "../../components/Card";

const DownloadPDF: React.FC = () => {
    return (
        <div className="p-6">
            <Card>
                <h1 className="text-xl font-bold mb-4">Download PDF</h1>
                <p>Download your result sheet as PDF.</p>
            </Card>
        </div>
    );
};

export default DownloadPDF;
