import React from 'react';
import './CustomPrompt.scss';

interface CustomPromptProps {
    value: string;
    onChange: (value: string) => void;
    className?: string;
}

const CustomPrompt: React.FC<CustomPromptProps> = ({ value, onChange, className = '' }) => {
    return (
        <input
            type="text"
            className={`custom-prompt ${className}`}
            value={value}
            onChange={(e) => onChange(e.target.value)}
            placeholder="Enter custom prompt"
        />
    );
};

export default CustomPrompt;