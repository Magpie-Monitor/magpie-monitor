import React from 'react';
import './LinkComponent.scss';

interface LinkComponentProps {
    href: string;
    children: React.ReactNode;
    className?: string;
    onClick?: React.MouseEventHandler<HTMLAnchorElement>;
}

const LinkComponent: React.FC<LinkComponentProps> = ({ href, children, className = '', onClick }) => {
    return (
        <a href={href} className={`link-component ${className}`} onClick={onClick}>
            {children}
        </a>
    );
};

export default LinkComponent;
