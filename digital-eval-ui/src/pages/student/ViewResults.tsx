import React from "react";
import Card from "../../components/Card";

const ViewResults: React.FC = () => {
    return (
        <div className="p-6">
            <Card>
                <h1 className="text-xl font-bold mb-4">View Results</h1>
                <p>Your semester results will be displayed here.</p>
            </Card>
        </div>
    );
};

export default ViewResults;
