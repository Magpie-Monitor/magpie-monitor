import React from 'react';
import './SectionComponent.scss';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';

interface SectionComponentProps {
    icon: string;
    title: React.ReactNode;
    children: React.ReactNode;
    callback?: () => void;
}

const SectionComponent: React.FC<SectionComponentProps> = ({ icon, title, children, callback }) => {
    return (
        <div className="section">
            <div className="section__header">
                <div className="section__icon">
                    <SVGIcon iconName={icon} />
                </div>
                <div className="section__title">{title}</div>
                {callback && (
                    <button
                        className="section__button"
                        onClick={callback}
                        aria-label="Add"
                    >
                        <SVGIcon iconName="plus-icon" />
                    </button>
                )}
            </div>
            <div className="section__divider"></div>
            <div className="section__content">{children}</div>
        </div>
    );
};

export default SectionComponent;
