import React from "react";

type Props = {
  children: React.ReactNode;
  className?: string;
  noPadding?: boolean;
  onClick?: (e: React.MouseEvent<HTMLDivElement>) => void;
};

const Card: React.FC<Props> = ({ children, className, noPadding, onClick }) => {
  return (
    <div
      onClick={onClick}
      className={`bg-white rounded-xl border border-slate-200 shadow-sm hover:shadow-md transition-all duration-300 ${
        noPadding ? "" : "p-6"
      } ${className ?? ""}`}
    >
      {children}
    </div>
  );
};

export default Card;