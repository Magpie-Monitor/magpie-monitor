import React from 'react';
import './CustomTag.scss';

interface CustomTagProps {
    logo?: React.ReactNode;
    name: string;
    onClick?: () => void;
    maxWidth?: string;
    maxHeight?: string;
}

const CustomTag =
    ({ logo, name, onClick, maxWidth = '250px', maxHeight = '40px' }: CustomTagProps) => {
    return (
        <div className="custom-tag" style={{ maxWidth, maxHeight }} onClick={onClick}>
            {logo && (
                <div className="custom-tag__logo">
                    {logo}
                </div>
            )}
            <div className="custom-tag__name">
                {name}
            </div>
        </div>
    );
};

export default CustomTag;
