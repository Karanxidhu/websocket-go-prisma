-- AlterTable
ALTER TABLE "MediaFile" ADD COLUMN     "message" TEXT,
ALTER COLUMN "url" DROP NOT NULL;
