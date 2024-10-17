import React from 'react';
import './PlaceholderComponent.scss';

interface PlaceholderComponentProps {
    icon?: React.ReactNode;
    title: React.ReactNode;
    children: React.ReactNode;
}

const PlaceholderComponent: React.FC<PlaceholderComponentProps> = ({icon, title, children}) => {
    return (
        <div className="placeholder-component">
            <div className="header">
                {icon && <div className="icon-container">{icon}</div>}
                <div>{title}</div>
            </div>
            <div className="divider"></div>
            <div className="payload">
                {children}
            </div>
        </div>
    );
};

export default PlaceholderComponent;
