import React from 'react';
import './SectionComponent.scss';
import SVGIcon from '@/components/SVGIcon/SVGIcon.tsx';

interface SectionComponentProps {
    icon: React.ReactNode;
    title: React.ReactNode;
    children: React.ReactNode;
}

const SectionComponent: React.FC<SectionComponentProps> = ({icon, title, children}) => {
    return (
        <div className="section">
            <div className="section__header">
                <div className="section__icon">
                <SVGIcon iconName={icon}/>
                </div>
                <div className="section__title">{title}</div>
            </div>
            <div className="section__divider"></div>
            <div className="section__content">
                {children}
            </div>
        </div>
    );
};

export default SectionComponent;
