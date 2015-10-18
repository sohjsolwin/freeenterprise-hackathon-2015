using System;
using System.Collections.Generic;
using System.Linq;
using System.Web;

namespace HackSite.Dto
{
    public class NewsData
    {
        public string State { get; set; }
        public string StateAbbr { get; set; }
        public string ArticleTitle { get; set; }
        public decimal? score { get; set; }
        public string ArticleLink { get; set; }
    }
}