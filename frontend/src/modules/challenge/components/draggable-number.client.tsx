//frontend/src/modules/challenge/components/draggable-number.client.tsx
"use client";

import { useDraggable } from "@dnd-kit/core";
import { CSS } from "@dnd-kit/utilities";
import { GAME } from "@/shared/config/gameStyles";

interface DraggableProps {
  id: string;
  value: string;
  onClick?: () => void;
  classNameExtra?: string;
}

export function DraggableNumber({
  id,
  value,
  onClick,
  classNameExtra = "",
}: DraggableProps) {
  const { attributes, listeners, setNodeRef, transform, isDragging } =
    useDraggable({ id });

  const style = {
    transform: CSS.Translate.toString(transform),
    transition: transform ? "none" : "transform 200ms ease",
    zIndex: isDragging ? 999 : 1,
    touchAction: "none",
    WebkitUserDrag: "none" as const,
  };

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...listeners}
      {...attributes}
      onClick={onClick}
      className={`${GAME.number(isDragging)} ${classNameExtra}`}>
      {value}
    </div>
  );
}
