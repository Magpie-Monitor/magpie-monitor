import React from 'react';
import './LinkComponent.scss';
import StateBadge from 'components/StateBadge/StateBadge';
import { useNavigate } from 'react-router-dom';

interface LinkComponentProps {
    to: string;
    children: React.ReactNode;
    className?: string;
    onClick?: React.MouseEventHandler<HTMLDivElement>;
    isRunning?: boolean;
}

const LinkComponent: React.FC<LinkComponentProps> = ({
    to,
    children,
    className = '',
    onClick,
    isRunning,
}) => {
    const navigate = useNavigate();

    const handleClick = (e: React.MouseEvent<HTMLDivElement>) => {
        if (onClick) {
            onClick(e);
        }

        navigate(to);
    };
    return (
        <div className="link-container">
            <div
                className={`link-component ${className} 
            link-component__${isRunning === false ? 'down' : 'up'}`}
                onClick={handleClick}
            >
                {children}
            </div>
            {isRunning !== undefined && (
                <div className="link-container__state-badge">
                    <StateBadge label={isRunning ? 'UP' : 'DOWN'} />
                </div>
            )}
        </div>
    );

};

export default LinkComponent;
