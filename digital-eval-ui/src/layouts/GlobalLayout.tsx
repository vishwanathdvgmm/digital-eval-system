import React from "react";
import Navbar from "../components/Navbar";
import Sidebar from "../components/Sidebar";
import TopbarUser from "../components/navigation/TopbarUser";

type Props = {
  children: React.ReactNode;
};

const GlobalLayout: React.FC<Props> = ({ children }) => {
  return (
    <div className="min-h-screen flex flex-col bg-gray-50">
      <Navbar>
        <div className="ml-auto">
          <TopbarUser />
        </div>
      </Navbar>

      <div className="flex flex-1">
        <aside className="w-64 hidden md:block border-r bg-white">
          <Sidebar />
        </aside>

        <main className="flex-1 p-6">
          <div className="max-w-7xl mx-auto">{children}</div>
        </main>
      </div>
    </div>
  );
};

export default GlobalLayout;