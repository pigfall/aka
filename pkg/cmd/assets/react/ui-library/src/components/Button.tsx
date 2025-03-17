import React, { MouseEventHandler } from 'react';

interface ButtonProps {
  children: React.ReactNode;
  onClick?: MouseEventHandler<HTMLButtonElement>;
  style?: React.CSSProperties;
}

const Button: React.FC<ButtonProps> = ({ children, onClick, style }) => {
  const defaultStyle: React.CSSProperties = {
    padding: '10px 20px',
    backgroundColor: '#007bff',
    color: 'white',
    border: 'none',
    borderRadius: '5px',
    cursor: 'pointer',
  };

  const combinedStyle: React.CSSProperties = { ...defaultStyle, ...style };

  return (
    <button style={combinedStyle} onClick={onClick}>
      {children}
    </button>
  );
};

export default Button;
