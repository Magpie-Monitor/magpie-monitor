import React, { useEffect, useRef } from 'react';
import './OverlayComponent.scss';

interface OverlayComponentProps {
    isDisplayed: boolean;
    onClose: () => void;
    children?: React.ReactNode;
}

const OverlayComponent: React.FC<OverlayComponentProps> = ({ isDisplayed, onClose, children }) => {
    const overlayRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (overlayRef.current && !overlayRef.current.contains(event.target as Node)) {
                onClose();
            }
        };

        const handleEscape = (event: KeyboardEvent) => {
            if (event.key === 'Escape') {
                onClose();
            }
        };

        if (isDisplayed) {
            document.addEventListener('mousedown', handleClickOutside);
            document.addEventListener('keydown', handleEscape);
        }

        return () => {
            document.removeEventListener('mousedown', handleClickOutside);
            document.removeEventListener('keydown', handleEscape);
        };
    }, [isDisplayed, onClose]);

    useEffect(() => {
        if (isDisplayed && overlayRef.current) {
            overlayRef.current.focus();
        }
    }, [isDisplayed]);

    return (
        isDisplayed && (
            <div className="modal-overlay">
                <div ref={overlayRef} className="modal-content">
                    {children || <p>List placeholder</p>}
                </div>
            </div>
        )
    );
};

export default OverlayComponent;
