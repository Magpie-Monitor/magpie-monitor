import { createRef, useEffect } from "react";
import "./Popup.scss";

interface PopupProps {
  isDisplayed: boolean;
  setIsDisplayed: (arg: boolean) => void;
  children: React.ReactNode;
}

export const Popup = ({
  isDisplayed,
  setIsDisplayed,
  children,
}: PopupProps) => {
  const popupRef = createRef<HTMLDivElement>();

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        popupRef.current &&
        !popupRef.current.contains(event.target as Node)
      ) {
        setIsDisplayed(false);
      }
    };
    const handleEscape = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        setIsDisplayed(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    document.addEventListener("keydown", handleEscape);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
      document.removeEventListener("keydown", handleEscape);
    };
  });

  useEffect(() => {
    if (isDisplayed) {
      popupRef.current?.focus();
    }
  }, [isDisplayed, popupRef]);

  return (
    isDisplayed && (
      <div ref={popupRef} className="popup">
        {children}
      </div>
    )
  );
};

export default Popup;
