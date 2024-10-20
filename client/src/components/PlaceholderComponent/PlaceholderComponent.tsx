import React from 'react';
import './PlaceholderComponent.scss';

interface PlaceholderComponentProps {
    icon: React.ReactNode;
    title: React.ReactNode;
    children: React.ReactNode;
}

const PlaceholderComponent: React.FC<PlaceholderComponentProps> = ({icon, title, children}) => {
    return (
        <div className="placeholder">
            <div className="placeholder__header">
                {icon && <div className="placeholder__icon">{icon}</div>}
                <div className="placeholder__title">{title}</div>
            </div>
            <div className="placeholder__divider"></div>
            <div className="placeholder__content">
                {children}
            </div>
        </div>
    );
};

export default PlaceholderComponent;
