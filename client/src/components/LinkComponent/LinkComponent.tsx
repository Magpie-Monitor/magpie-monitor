import React from 'react';
import './LinkComponent.scss';

interface LinkComponentProps {
    href: string;
    children: React.ReactNode;
    className?: string;
}

const LinkComponent: React.FC<LinkComponentProps> = ({ href, children, className = '' }) => {
    return (
        <a href={href} className={`link-component ${className}`}>
            {children}
        </a>
    );
};

export default LinkComponent;
