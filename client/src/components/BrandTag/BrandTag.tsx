import './BrandTag.scss';

interface BrandTagProps {
    logo: React.ReactNode;
    name: string;
}

const BrandTag = ({ logo, name }: BrandTagProps) => {
    return (
        <div className="brand-tag">
            <div className="brand-tag__logo">{logo}</div>
            <div className="brand-tag__name">{name}</div>
        </div>
    );
};

export default BrandTag;