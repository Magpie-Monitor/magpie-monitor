import React from 'react';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';
import './DeleteIconButton.scss';

interface DeleteIconButtonProps {
    onClick: () => void;
}

const DeleteIconButton: React.FC<DeleteIconButtonProps> = ({ onClick }) => {
    return (
        <button className="delete-icon-button" onClick={onClick} aria-label="Delete">
            <SVGIcon iconName="delete-icon" />
        </button>
    );
};

export default DeleteIconButton;
