import React from 'react';
import './Spinner.scss';

interface SpinnerProps {
    size?: string;
}

const Spinner: React.FC<SpinnerProps> = ({ size = '40px' }) => {
    return (
      <div
        className="spinner"
        style={{ width: size, height: size }}
      />
    );
};

export default Spinner;
