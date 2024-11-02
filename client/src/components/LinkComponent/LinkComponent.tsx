import React from 'react';
import './LinkComponent.scss';
import StateBadge from 'components/StateBadge/StateBadge';

interface LinkComponentProps {
    href: string;
    children: React.ReactNode;
    className?: string;
    onClick?: React.MouseEventHandler<HTMLAnchorElement>;
    isRunning?: boolean;
}

const LinkComponent: React.FC<LinkComponentProps> = ({ href, children, className = '', onClick, isRunning }) => {
    return (
        <div className="link-container">
            <a href={href} className={`link-component ${className} link-component__${isRunning === false ? 'down' : 'up'}`} onClick={onClick}>
                {children}
            </a>
            {isRunning !== undefined && (
                <StateBadge label={isRunning ? 'UP' : 'DOWN'} />
            )}
        </div>
    );
};

export default LinkComponent;
